package state

import (
	// "crypto/sha256"
	// "crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"

	"client/keyvalue"
	"client/settings"
	"client/thrift/pool"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Recover struct {
	File *keyvalue.File
}

func (r Recover) Run(sm *StateMachine) {
	file := r.File
	resp, err := util.Post(sm.Options, fmt.Sprintf("file/%s/download", file.Id), map[string]string{"clientId": sm.Options.ClientId})
	if err != nil {
		log.Println("Unable to connect to server to recover a file: ", err)
		return
	}
	if resp["error"] != nil {
		log.Println("Error recovering a file ", resp["error"], " ", resp["message"])
		return
	}
	if resp["allowed"].(bool) {
		file.EncodedSize = int64(resp["size"].(float64))
		file.UnencodedSize = int64(resp["originalSize"].(float64))

		blockIds := resp["blocks"].([]interface{})
		for _, bId := range blockIds {
			blockId := bId.(string)
			file.Blocks = append(file.Blocks, keyvalue.Block{Id: blockId, Shards: make([]keyvalue.Shard, settings.M)})
		}

		clients := resp["clients"].([]interface{})
		for _, c := range clients {
			client := c.(map[string]interface{})
			blockId := client["blockId"].(string)
			offset := int(client["offset"].(float64))
			for _, block := range file.Blocks {
				if block.Id == blockId {
					block.Shards[offset].Id = client["id"].(string)
					block.Shards[offset].IP = client["IP"].(string)
					block.Shards[offset].Hash, err = base64.StdEncoding.DecodeString(client["hash"].(string))
					if err != nil {
						log.Println("Couldn't decode shard hash ", err)
					}
					break
				}
			}
		}

		clientPool := pool.NewClientPool(thrift.NewTBufferedTransportFactory(10000), thrift.NewTBinaryProtocolFactoryDefault())
		defer clientPool.CloseAll()

		decode := &Decode{
			File:          file,
			EncodedBlocks: make([][]byte, len(file.Blocks)),
			BlockErrs:     make([][]byte, len(file.Blocks)),
		}

		for b, block := range file.Blocks {
			for _, shard := range block.Shards {
				client, err := clientPool.GetClient(shard.IP)
				if err != nil {
					log.Println("Error opening connection to ", shard.IP, err)
					continue
				}

				replica, err := client.Download(shard.Id)
				if err != nil {
					log.Println("Error during download:", err)
					decode.BlockErrs[b] = append(decode.BlockErrs[b], byte(b))
				} else {
					// TODO: Fix this so shard hashes work
					// shardHash := sha256.New()
					// shardHash.Write(replica.Shard)
					// log.Println("Hash: ", string(shardHash.Sum(nil)), string(shard.Hash))
					// if subtle.ConstantTimeCompare(shardHash.Sum(nil), shard.Hash) == 1 {
					decode.EncodedBlocks[b] = append(decode.EncodedBlocks[b], replica.Shard...)
					// } else {
					// 	decode.BlockErrs[b] = append(decode.BlockErrs[b], byte(b))
					// }
				}
			}
		}
		sm.Add(decode)
	}
}
