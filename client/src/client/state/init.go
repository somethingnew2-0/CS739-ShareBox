package state

import (
	"fmt"
	"net/http"
	"net/url"

	"client/settings"
)

type Init struct{}

func (i Init) Run(sm *StateMachine) {

	sm.Options.LoadFromJSON()
	if sm.Options.ClientId == "" || sm.Options.UserId == "" {
		http.PostForm(fmt.Sprintf("%s/user/new", settings.ServerAddress), url.Values{})

	}

	// TODO: Connect to consul and report status

	// TODO: Recover files
	// TODO: Upload new files

	// Say we have 1.5 GB of space available
	// http.PostForm(fmt.Sprintf("%s/client/%s/init", settings.ServerAddress, sm.ClientId), url.Values{"space": {"1610612736"}})

	sm.Add(Watch{})
}
