package state

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"client/keyvalue"
	"client/util"
)

type Init struct{}

func (i Init) Run(sm *StateMachine) {

	sm.Options.Load()
	if sm.Options.ClientId == "" || sm.Options.UserId == "" {
		sm.Options.HashPassword()
		resp, err := util.Post(sm.Options, "user/new", map[string]string{"username": sm.Options.Username, "passwordHash": string(sm.Options.Hash)})
		if err != nil {
			log.Fatal("Couldn't connect and create new user with server: ", err)
		}
		userResp := resp["user"].(map[string]interface{})
		sm.Options.UserId = userResp["id"].(string)
		sm.Options.ClientId = userResp["clientId"].(string)
		sm.Options.AuthToken = userResp["auth"].(string)

		sm.Options.Save()

		// TODO: Actually get user's available disk space (instead of just 1GB)
		resp, err = util.Post(sm.Options, fmt.Sprintf("client/%s/init", sm.Options.ClientId), map[string]interface{}{"space": 1 << 30})
		if err != nil {
			log.Fatal("Couldn't init client with server: ", err)
		}
	}

	resp, err := util.Get(sm.Options, fmt.Sprintf("client/%s/status", sm.Options.ClientId))
	if err != nil {
		log.Fatal("Couldn't get status of client from server: ", err)
	}
	fresh := resp["new"].(bool)
	recovery := resp["recovery"].(bool)

	if fresh {
		err = filepath.Walk(sm.Options.Dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				sm.Add(&Read{Create: true, Path: path, Info: info})
			}
			return nil
		})
		if err != nil {
			log.Fatal("Error when walking data directory", err)
		}

	} else if recovery {
		resp, err = util.Get(sm.Options, fmt.Sprintf("client/%s/recover", sm.Options.ClientId))
		if err != nil {
			log.Fatal("Unable to connect to server to recover files: ", err)
		}
		if resp["allowed"].(bool) {
			files := resp["fileList"].([]map[string]string)
			for _, file := range files {
				sm.Add(Recover{File: &keyvalue.File{Id: file["id"], Name: file["name"], Hash: []byte(file["hash"])}})
			}
		}
	}
	sm.Add(Watch{})
	sm.Add(Health{})
	sm.Add(Replica{})
}
