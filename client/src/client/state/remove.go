package state

import (
	"fmt"
	"log"

	"client/thrift/pool"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Remove struct {
	Path string
}

func (r Remove) Run(sm *StateMachine) {
	file, err := sm.Files.GetFile(r.Path)
	var resp map[string]interface{}
	if err == nil {
		resp, err = util.Post(sm.Options, fmt.Sprintf("client/%s/file/remove", sm.Options.ClientId), map[string]string{"id": file.Id, "name": r.Path})
		if err != nil {
			log.Println("Error removing file", err)
			return
		}
	} else {
		resp, err = util.Post(sm.Options, fmt.Sprintf("client/%s/file/remove", sm.Options.ClientId), map[string]string{"name": r.Path})
		if err != nil {
			log.Println("Error removing file", err)
			return
		}
	}

	clientPool := pool.NewClientPool(thrift.NewTBufferedTransportFactory(10000), thrift.NewTBinaryProtocolFactoryDefault())
	defer clientPool.CloseAll()
	if resp["allowed"].(bool) {
		shards := resp["clients"].([]interface{})
		for _, s := range shards {
			shard := s.(map[string]interface{})
			// TODO Validate this an IP address using net.IP
			client, err := clientPool.GetClient(shard["IP"].(string))
			if err != nil {
				log.Println("Error opening connection to ", shard["IP"].(string), err)
				continue
			}

			err = client.Remove(shard["id"].(string))
			if err != nil {
				log.Println("Error during remove", err)
			}
		}

		if file != nil {
			resp, err = util.Post(sm.Options, fmt.Sprintf("file/%s/delete", file.Id), map[string]string{"clientId": sm.Options.ClientId})
			if err != nil {
				log.Println("Error deleting file", err)
				return
			}
			if !resp["success"].(bool) {
				log.Println("Deleting file from server was unsucessful")
			}
		}
	}
}
