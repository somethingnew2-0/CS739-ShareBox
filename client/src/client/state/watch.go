package state

import (
	"log"
	"os"

	"gopkg.in/fsnotify.v1"
)

type Watch struct{}

func (w Watch) Run(sm *StateMachine) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(sm.Options.Dir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println("Watch event: ", event.Op, event.Name)
			switch event.Op {
			case fsnotify.Create:
				info, err := os.Stat(event.Name)
				if err != nil {
					log.Println("Error stating file: ", event.Name, err)
				}
				sm.Add(&Create{Path: event.Name, Info: info})
			case fsnotify.Write:
				// TODO: This could be a race condition here
				sm.Add(&Remove{Path: event.Name})
				sm.Add(&Create{Path: event.Name})
			case fsnotify.Remove:
				sm.Add(&Remove{Path: event.Name})
			case fsnotify.Rename:
				// TODO: This could be a race condition here
				sm.Add(&Remove{Path: event.Name})
				sm.Add(&Create{Path: event.Name})
				log.Println("Rename file: ", event.Name)
			}
		case err := <-watcher.Errors:
			log.Println("Error watching directory:", err)
		}
	}
}
