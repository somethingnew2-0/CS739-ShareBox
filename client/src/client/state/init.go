package state

import (
	"fmt"
	"net/http"
	"net/url"

	"client/settings"
)

type Init struct{}

func (i Init) Run(sm *StateMachine) {

	if sm.Options.ClientId == nil || sm.Options.UserId == nil {
		if _, err := os.Stat(settings.ConfigFile); err == nil {

		} else {
			http.PostForm(fmt.Sprintf("%s/user/new", settings.ServerAddress), url.Values{})

		}
	}
	// Say we have 1.5 GB of space available
	// http.PostForm(fmt.Sprintf("%s/client/%s/init", settings.ServerAddress, sm.ClientId), url.Values{"space": {"1610612736"}})

	sm.Add(Watch{})
}
