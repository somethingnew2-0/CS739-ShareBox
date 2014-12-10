package state

import (
	"fmt"
	"log"

	"client/keyvalue"
	"client/settings"
	"client/thrift/pool"
	"client/thrift/replica"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Upload struct {
	EncodedBlocks [][]byte
	File          *keyvalue.File
}

func (u Upload) Run(sm *StateMachine) {
	resp, err := util.Post(sm.Options, fmt.Sprintf("client/%s/file/add", sm.Options.ClientId), u.File)
	if err != nil {
		log.Println("Error adding file", err)
		return
	}

	if resp["allowed"].(bool) {
		file := u.File
		file.Id = resp["id"].(string)
		blockIds := resp["blocks"].([]interface{})
		for i, blockId := range blockIds {
			file.Blocks[i].Id = blockId.(string)
		}

		clients := resp["clients"].([]interface{})
		for _, c := range clients {
			client := c.(map[string]interface{})
			blockId := client["blockId"]
			offset := client["offset"].(int)
			for _, block := range file.Blocks {
				if block.Id == blockId {
					shard := block.Shards[offset]
					shard.Id = client["id"].(string)
					// TODO Validate this an IP address using net.IP
					shard.IP = client["IP"].(string)
					break
				}
			}
		}

		transportPool := pool.NewTransportPool(thrift.NewTBufferedTransportFactory(10000))
		defer transportPool.CloseAll()
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

		for b, block := range file.Blocks {
			for s, shard := range block.Shards {
				transport, err := transportPool.GetTransport(shard.IP)
				if err != nil {
					log.Println("Error opening connection to ", shard.IP, err)
					continue
				}

				client := replica.NewReplicatorClientFactory(transport, protocolFactory)
				client.Ping()

				shardData := u.EncodedBlocks[b][s*settings.ShardLength : (s+1)*settings.ShardLength]
				err = client.Add(&replica.Replica{
					Shard:       shardData,
					ShardHash:   shard.Hash,
					ShardOffset: int32(s),
					ShardId:     shard.Id,
					BlockId:     block.Id,
					FileId:      file.Id,
					ClientId:    sm.Options.ClientId,
				})

				if err != nil {
					log.Println("Error during upload", err)
				}
			}
		}

		resp, err := util.Post(sm.Options, fmt.Sprintf("file/%s/commit", file.Id), map[string]string{"clientId": sm.Options.ClientId})
		if err != nil {
			log.Println("Error commiting file", err)
			return
		}

		if resp["success"].(bool) {
			sm.Files.SetFile(file.Name, file)
		} else {
			log.Println("Commiting file was unsuccessful", err)

		}
	} else {
		log.Println("File upload not allowed")
		return
	}
}
