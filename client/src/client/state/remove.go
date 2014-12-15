package state

import (
	"fmt"
	"log"

	"client/keyvalue"
	"client/util"
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
		if resp["error"] != nil {
			log.Println("Error removing file ", resp["error"], " ", resp["message"])
			return
		}
	} else {
		resp, err = util.Post(sm.Options, fmt.Sprintf("client/%s/file/remove", sm.Options.ClientId), map[string]string{"name": r.Path})
		if err != nil {
			log.Println("Error removing file", err)
			return
		}
		if resp["error"] != nil {
			log.Println("Error removing file ", resp["error"], " ", resp["message"])
			return
		}
	}

	if resp["allowed"].(bool) {
		invalidate := make([]keyvalue.Shard, 0)
		shards := resp["clients"].([]interface{})
		for _, s := range shards {
			shard := s.(map[string]interface{})
			// TODO Validate this an IP address using net.IP
			invalidate = append(invalidate, keyvalue.Shard{
				Id: shard["id"].(string),
				IP: shard["IP"].(string),
			})
		}

		sm.Add(&Invalidate{
			Shards: invalidate,
			File:   file,
			CallBack: func(i *Invalidate, sm *StateMachine) {
				if i.File != nil {
					resp, err = util.Post(sm.Options, fmt.Sprintf("file/%s/delete", i.File.Id), map[string]string{"clientId": sm.Options.ClientId})
					if err != nil {
						log.Println("Error deleting file", err)
						return
					}
					if resp["error"] != nil {
						log.Println("Error deleting file ", resp["error"], " ", resp["message"])
						return
					}
					if !resp["success"].(bool) {
						log.Println("Deleting file from server was unsucessful")
					}
				}
			},
		})
	}
}
