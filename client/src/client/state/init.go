package state

import (
	"fmt"
	"net/http"
	"net/url"

	"client/settings"
)

type Init struct{}

func (i Init) Run(sm *StateMachine) {
	// Say we have 1.5 GB of space available
	http.PostForm(fmt.Sprintf("%s/client/%s/init", settings.ServerAddress, sm.ClientId),
		url.Values{"space": {"1610612736"}})
}
