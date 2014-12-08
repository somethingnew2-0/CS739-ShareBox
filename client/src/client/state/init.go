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
		resp, err := util.Post(fmt.Sprintf("%s/user/new", settings.ServerAddress), url.Values{})
		if err != nil {
			log.Fatal("Couldn't connect and create new user with server", err)
		}
		userResp := resp["user"].(map[string]interface{})
		sm.Options.UserId = userResp["id"].(string)
		sm.Options.ClientId = userResp["clientId"].(string)

		sm.Options.Save()

		// TODO: Actually get user's available disk space (instead of just 1GB)
		resp, err = util.Post(fmt.Sprintf("%s/client/%s/init", settings.ServerAddress), url.Values{"space": {string(1 << 30)}})
		if err != nil {
			log.Fatal("Couldn't init client with server", err)
		}

		err = filepath.Walk(sm.Options.Dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				encrypt := &Encrypt{}

				blocks := math.Ceil(float64(info.Size()) / float64(settings.BlockSize))
				// TODO: Do I add the file hash here?
				encrypt.File = &keyvalue.File{Name: path, Size: info.Size()}
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
		// TODO: Upload new files
	} else {

		// TODO: Recover files
	}

	// TODO: Connect to consul and report status

	// Say we have 1.5 GB of space available
	// http.PostForm(fmt.Sprintf("%s/client/%s/init", settings.ServerAddress, sm.ClientId), url.Values{"space": {"1610612736"}})

	sm.Add(Watch{})
}
