package main

import (
	"client/settings"
	"client/state"
)

func main() {
	stateMachine := state.NewStateMachine(settings.ClientId)

	stateMachine.Add(state.Init{})
	stateMachine.Run()
}
