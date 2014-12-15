package state

import (
	"fmt"
	"log"

	"client/keyvalue"
	"client/util"
)

type Create struct {
	EncodedBlocks [][]byte
	File          *keyvalue.File
}

func (u Create) Run(sm *StateMachine) {
	resp, err := util.Post(sm.Options, fmt.Sprintf("client/%s/file/add", sm.Options.ClientId), u.File)
	if err != nil {
		log.Println("Error adding file", err)
		return
	}
	if resp["error"] != nil {
		log.Println("Error adding file ", resp["error"], " ", resp["message"])
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
		sm.Add(&Upload{
			EncodedBlocks: u.EncodedBlocks,
			File:          file,
		})
	} else {
		log.Println("File upload not allowed", u.File.Name)
		return
	}
}
