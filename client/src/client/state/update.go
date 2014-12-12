package state

import (
	"crypto/subtle"
	"fmt"
	"log"
	"time"

	"client/keyvalue"
	"client/settings"
	"client/thrift/pool"
	"client/thrift/replica"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Update struct {
	EncodedBlocks [][]byte
	File          *keyvalue.File
}

func (u Update) Run(sm *StateMachine) {
	f, err := sm.Files.GetFile(u.File.Id)

	// Updates usually follow a Create.  Wait a bit for it to complete.
	// If the Create never happens or takes too long to finish, we'll do it ourselves.
	failures := 0
	ticker := time.NewTicker(time.Second)
	for err != nil {
		// If we reach the timeout, just create the file instead
		if failures > settings.UpdateTimeout {
			sm.Add(&Create{EncodedBlocks: u.EncodedBlocks, File: u.File})
			return
		}
		<-ticker.C
		f, err = sm.Files.GetFile(u.File.Id)
	}

	if subtle.ConstantTimeCompare([]byte(u.File.Hash), []byte(f.Hash)) == 1 {
		log.Println("No changes were actaully detected with file update")
		return
	}

	u.File.Id = f.Id

	for i, block := range u.File.Blocks {
		if len(f.Blocks) < i {
			u.File.Blocks[i].Id = f.Blocks[i].Id
			if subtle.ConstantTimeCompare([]byte(block.Hash), []byte(f.Blocks[i].Hash)) == 1 {
				// The block hasn't changed
				// TODO: Delete this block from the posted JSON becuase it hasn't changed
				continue
			}
			// The block has changed
		}
	}

	// Check if the file shrunk
	if len(u.File.Blocks) < len(f.Blocks) {
		// TODO: Figure out what to do here
		// Call block delete for discarded blocks?
	}

	resp, err := util.Post(sm.Options, fmt.Sprintf("client/%s/file/update", sm.Options.ClientId), u.File)
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
			offset := int(client["offset"].(float64))
			for _, block := range file.Blocks {
				if block.Id == blockId {
					block.Shards[offset].Id = client["id"].(string)
					// TODO Validate this an IP address using net.IP
					block.Shards[offset].IP = client["IP"].(string)
					break
				}
			}
		}

		clientPool := pool.NewClientPool(thrift.NewTBufferedTransportFactory(10000), thrift.NewTBinaryProtocolFactoryDefault())
		defer clientPool.CloseAll()

		for b, block := range file.Blocks {
			for s, shard := range block.Shards {
				client, err := clientPool.GetClient(shard.IP)
				if err != nil {
					log.Println("Error opening connection to ", shard.IP, err)
					continue
				}

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
