package state

import (
	"fmt"
	"log"
	"strconv"

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
	resp, err := util.Post(fmt.Sprintf("client/%s/file/add", sm.Options.ClientId), u.File)
	if err != nil {
		log.Println("Error adding file", err)
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
					// TODO Validate this an IP address using net.IP
					shard.IP = client["IP"]
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

		resp, err := util.Post(fmt.Sprintf("file/%s/commit", file.Id), map[string]string{"clientId": sm.Options.ClientId})
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
