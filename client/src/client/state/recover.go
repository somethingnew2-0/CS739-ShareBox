package state

import (
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"client/keyvalue"
	"client/settings"
	"client/thrift/replica"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Recover struct {
	File *keyvalue.File
}

func (r Recover) Run(sm *StateMachine) {
	file := r.File
	resp, err := util.Post(fmt.Sprintf("file/%s/download", file.Id), url.Values{"clientId": {sm.Options.ClientId}})
	if err != nil {
		log.Println("Unable to connect to server to recover a file: ", err)
		return
	}
	if allowed, err := strconv.ParseBool(resp["allowed"].(string)); err == nil && allowed {
		file.Size, err = strconv.ParseInt(resp["size"].(string), 10, 64)
		if err != nil {
			log.Println("Cannot parse file size", err)
			return
		}

		blockIds := resp["blocks"].([]string)
		for _, blockId := range blockIds {
			file.Blocks = append(file.Blocks, keyvalue.Block{Id: blockId, Shards: make([]keyvalue.Shard, settings.M)})
		}

		clients := resp["clients"].([]map[string]string)
		for _, client := range clients {
			blockId := client["blockId"]
			offset, err := strconv.ParseInt(client["offset"], 10, 64)
			if err != nil {
				log.Println("Cannot parse shard offset", err)
				break
			}
			for _, block := range file.Blocks {
				if block.Id == blockId {
					shard := block.Shards[offset]
					shard.Id = client["id"]
					shard.Hash = client["Hash"]
					shard.IP = client["IP"]
					break
				}
			}
		}

		transports := make(map[string]thrift.TTransport)
		transportFactory := thrift.NewTBufferedTransportFactory(10000)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

		decode := &Decode{
			File:          file,
			EncodedBlocks: make([][]byte, len(file.Blocks)),
			BlockErrs:     make([][]byte, len(file.Blocks)),
		}

		for b, block := range file.Blocks {
			for _, shard := range block.Shards {
				transport := transports[shard.IP]
				if transport == nil {
					addr := fmt.Sprintf("%s:%d", shard.IP, settings.ClientPort)
					if settings.ClientTLS {
						cfg := new(tls.Config)
						cfg.InsecureSkipVerify = true
						transport, err = thrift.NewTSSLSocket(addr, cfg)
					} else {
						transport, err = thrift.NewTSocket(addr)
					}
					if err != nil {
						log.Println("Error opening socket:", err)
						break
					}
					transport = transportFactory.GetTransport(transport)
					defer transport.Close()
					if err := transport.Open(); err != nil {
						break
					}
					transports[shard.IP] = transport
				}
				client := replica.NewReplicatorClientFactory(transport, protocolFactory)
				client.Ping()

				replica, iv, err := client.Download(shard.Id)
				if iv != nil {
					log.Println("Invalid operation:", iv)
					decode.BlockErrs[b] = append(decode.BlockErrs[b], byte(b))
				} else if err != nil {
					log.Println("Error during download:", err)
					decode.BlockErrs[b] = append(decode.BlockErrs[b], byte(b))
				} else {
					shardHash := sha256.New()
					shardHash.Write(replica.Shard)
					if subtle.ConstantTimeCompare(shardHash.Sum(nil), []byte(shard.Hash)) == 1 {
						decode.EncodedBlocks[b] = append(decode.EncodedBlocks[b], replica.Shard...)
					} else {
						decode.BlockErrs[b] = append(decode.BlockErrs[b], byte(b))
					}
				}
			}
		}
		sm.Add(decode)
	}
}
