package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type Init struct{}

func (*Init) Run(sm *StateMachine) {
	// Say we have 1.5 GB of space available
	http.PostForm(fmt.Sprintf("%s/client/%s/init", ServerAddress, sm.ClientId),
		url.Values{"space": {"1610612736"}})
}
