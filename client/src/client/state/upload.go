package state

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"client/keyvalue"
	"client/settings"
	"client/util"
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

	resp, err := util.Post(fmt.Sprintf("%s/client/%s/file/add", settings.ServerAddress, sm.Options.ClientId), url.Values{"Request": {string(fileJson)}})
	if err != nil {
		log.Println("Error adding new file", err)
		return
	}

	if allowed, err := strconv.ParseBool(resp["allowed"].(string)); err == nil && allowed {
		file := u.File
		file.Id = resp["id"].(string)
		blockIds := resp["blocks"].([]string)
		for i, blockId := range blockIds {
			file.Blocks[i].Id = blockId
		}

		clients := resp["clients"].([]map[string]string)
		for _, client := range clients {
			blockId := client["blockId"]
			offset, err := strconv.ParseInt(client["offset"], 10, 32)
			if err != nil {
				log.Println("Cannot parse shard offset", err)
				break
			}

			for _, block := range file.Blocks {
				if block.Id == blockId {
					shard := block.Shards[offset]
					shard.Id = client["id"]
					shard.IP = net.ParseIP(client["IP"])
					break
				}
			}
		}
		sm.Files.SetFile(file.Name, file)
	} else {
		log.Println("File upload not allowed")
		return
	}
}
