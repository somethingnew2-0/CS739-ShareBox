package main

import (
	"log"
	"runtime"

	"client/settings"
	"client/state"

	"github.com/jessevdk/go-flags"
)

var (
	opts *settings.Options
	args []string
)

func init() {
	var err error
	args, err = flags.Parse(&opts)
	if err != nil {
		log.Fatalf("Error parsing options: &v\n", err)
	}

	// Set runtime GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	stateMachine := state.NewStateMachine(opts)

	stateMachine.Add(state.Init{})
	stateMachine.Run()
}
