package state

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

type Watch struct{}

func (w Watch) Run(sm *StateMachine) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(filepath.Join(os.TempDir(), "foo"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}
