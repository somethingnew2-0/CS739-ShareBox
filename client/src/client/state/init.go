package state

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"os"
	"path/filepath"

	"client/keyvalue"
	"client/settings"
	"client/util"
)

type Init struct{}

func (i Init) Run(sm *StateMachine) {

	sm.Options.Load()
	if sm.Options.ClientId == "" || sm.Options.UserId == "" {
		sm.Options.HashPassword()
		resp, err := util.Post("user/new", url.Values{"username": {}, "password": {string(sm.Options.Hash)}})
		if err != nil {
			log.Fatal("Couldn't connect and create new user with server: ", err)
		}
		userResp := resp["user"].(map[string]interface{})
		sm.Options.UserId = userResp["id"].(string)
		sm.Options.ClientId = userResp["clientId"].(string)

		sm.Options.Save()

		// TODO: Actually get user's available disk space (instead of just 1GB)
		resp, err = util.Post(fmt.Sprintf("client/%s/init", sm.Options.ClientId), url.Values{"space": {string(1 << 30)}})
		if err != nil {
			log.Fatal("Couldn't init client with server: ", err)
		}
	}

	resp, err := util.Get(fmt.Sprintf("client/%s/status", sm.Options.ClientId))
	if err != nil {
		log.Fatal("Couldn't get status of client from server: ", err)
	}
	fresh := resp["new"].(bool)
	recovery := resp["recovery"].(bool)

	if fresh {
		err = filepath.Walk(sm.Options.Dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				encrypt := &Encrypt{}

				blocks := math.Ceil(float64(info.Size()) / float64(settings.BlockSize))

				// Use encoded file size
				encrypt.File = &keyvalue.File{Name: path, EncodedSize: int64(blocks) * settings.BlockSize, UnencodedSize: info.Size()}
				if f, err := os.Open(path); err == nil {
					zeroBytes := (int64(blocks) * settings.BlockSize) - info.Size()
					data, _ := ioutil.ReadAll(f)
					encrypt.Plaintext = append(data, make([]byte, zeroBytes)...)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal("Error when walking data directory", err)
		}

	} else if recovery {
		resp, err = util.Get(fmt.Sprintf("client/%s/recover", sm.Options.ClientId))
		if err != nil {
			log.Fatal("Unable to connect to server to recover files: ", err)
		}
		if resp["allowed"].(bool) {
			files := resp["fileList"].([]map[string]string)
			for _, file := range files {
				sm.Add(Recover{File: &keyvalue.File{Id: file["id"], Name: file["name"], Hash: file["hash"]}})
			}
		}
		// TODO: Recover files

	} else {
		sm.Add(Watch{})
	}

	// TODO: Connect to consul and report status

}
