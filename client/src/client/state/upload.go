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
	clientPool := pool.NewClientPool(thrift.NewTBufferedTransportFactory(10000), thrift.NewTBinaryProtocolFactoryDefault())
	defer clientPool.CloseAll()

	for b, block := range u.File.Blocks {
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
				FileId:      u.File.Id,
				ClientId:    sm.Options.ClientId,
			})

			if err != nil {
				log.Println("Error during upload", err)
			}
		}
	}
	resp, err := util.Post(sm.Options, fmt.Sprintf("file/%s/commit", u.File.Id), map[string]string{"clientId": sm.Options.ClientId})
	if err != nil {
		log.Println("Error commiting file", err)
		return
	}

	if resp["success"].(bool) {
		sm.Files.SetFile(u.File.Name, u.File)
	} else {
		log.Println("Commiting file was unsuccessful", u.File.Name, err)

	}

}
