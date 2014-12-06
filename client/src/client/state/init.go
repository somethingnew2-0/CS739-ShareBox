package state

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"client/settings"
)

type Init struct{}

func (i Init) Run(sm *StateMachine) {

	sm.Options.Load()
	if sm.Options.ClientId == "" || sm.Options.UserId == "" {
		http.PostForm(fmt.Sprintf("%s/user/new", settings.ServerAddress), url.Values{})
		sm.Options.Save()

		err := filepath.Walk(sm.Options.Dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {

			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
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
