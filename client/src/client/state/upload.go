package state

import (
	"crypto/tls"
	"encoding/json"
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

type Upload struct {
	EncodedBlocks [][]byte
	File          *keyvalue.File
}

func (u Upload) Run(sm *StateMachine) {
	fileJson, err := json.Marshal(u.File)
	if err != nil {
		log.Println("Error json encoding file", err)
		return
	}

	resp, err := util.Post(fmt.Sprintf("client/%s/file/add", sm.Options.ClientId), url.Values{"Request": {string(fileJson)}})
	if err != nil {
		log.Println("Error adding new file", err)
		return
	}

	if resp["allowed"].(bool) {
		file := u.File
		file.Id = resp["id"].(string)
		blockIds := resp["blocks"].([]string)
		for i, blockId := range blockIds {
			file.Blocks[i].Id = blockId
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
					shard.IP = client["IP"]
					break
				}
			}
		}

		transports := make(map[string]thrift.TTransport)
		transportFactory := thrift.NewTBufferedTransportFactory(10000)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

		for b, block := range file.Blocks {
			for s, shard := range block.Shards {
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

				shardData := u.EncodedBlocks[b][s*settings.ShardLength : (s+1)*settings.ShardLength]
				iv, err := client.Add(&replica.Replica{
					Shard:       shardData,
					ShardHash:   shard.Hash,
					ShardOffset: int32(s),
					ShardId:     shard.Id,
					BlockId:     block.Id,
					FileId:      file.Id,
					ClientId:    sm.Options.ClientId,
				})

				if iv != nil {
					log.Println("Invalid operation:", iv)
				} else if err != nil {
					log.Println("Error during upload", err)
				}
			}
		}

		sm.Files.SetFile(file.Name, file)
	} else {
		log.Println("File upload not allowed")
		return
	}
}
