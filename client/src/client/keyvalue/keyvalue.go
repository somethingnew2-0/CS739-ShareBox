package keyvalue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxSetsPerSec uint = 1 << 15

type set struct {
	Key   string
	Value string
}

type KeyValue struct {
	LogDir         string
	store          map[string]string
	storeLock      *sync.RWMutex // Maps aren't thread safe, must lock on writes using a readers-writer lock
	pending        chan *set     // Pending sets are sent to channel to be added
	pendingPersist chan *set
}

func Init(logDir string) *KeyValue {
	log.Println("KeyValue starting at ", logDir)

	server := &KeyValue{
		LogDir:         logDir,
		store:          make(map[string]string),
		storeLock:      &sync.RWMutex{},
		pending:        make(chan *set, MaxSetsPerSec),
		pendingPersist: make(chan *set, MaxSetsPerSec),
	}

	os.MkdirAll(logDir, 0777)

	server.recover()
	log.Println("KeyValue fully recovered at", logDir)

	go server.set()

	go server.persistDelta()
	go server.persistBase()

	/*go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%v", server.store)
		})

		log.Fatal(http.ListenAndServe(":8080", nil))
	}()*/

	log.Println("KeyValue accepting requests at ", logDir)
	return server
}

func (s *KeyValue) recover() {
	entries, err := ioutil.ReadDir(s.LogDir)
	if err != nil {
		log.Printf("Error reading log directory, unable to recover: %v\n", err)
		return
	}

	names := make([]string, len(entries))
	for index, entry := range entries {
		names[index] = entry.Name()
	}
	sort.Strings(names)

	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	// Find the most recent back backup
	var baseEpoch int64
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if strings.LastIndex(name, "-base") >= 0 {
			split := strings.Split(name, "-")
			if len(split) == 2 {
				epoch, err := strconv.ParseInt(split[0], 10, 64)
				if err == nil {
					baseEpoch = epoch
				}
			}

			data, err := ioutil.ReadFile(path.Join(s.LogDir, name))
			if err != nil {
				log.Printf("Error reading base log, unable to recover: %v", err)
				return
			}

			err = json.Unmarshal(data, &s.store)
			if err != nil {
				log.Printf("Error unmarshalling base log, unable to recover: %v", err)
				return
			}

			// Truncate the list of names so we don't have to iterate
			// through all of them for delta recovery
			if len(names) > i+1 {
				names = names[i+1:]
			} else {
				// No further delta updates in the list
				return
			}

			break
		}
	}

	for _, name := range names {
		if strings.LastIndex(name, "-delta") >= 0 {
			split := strings.Split(name, "-")
			if len(split) == 2 {
				epoch, err := strconv.ParseInt(split[0], 10, 64)
				if err == nil && epoch > baseEpoch {
					data, err := ioutil.ReadFile(path.Join(s.LogDir, fmt.Sprintf("%d-delta", epoch)))
					if err != nil {
						log.Printf("Error reading delta log, recovery could be paritally incorrect: %v", err)
						continue
					}

					var sets []set
					err = json.Unmarshal(data, &sets)
					if err != nil {
						log.Printf("Error reading delta log, recovery could be paritally incorrect: %v", err)
						continue
					}

					for _, set := range sets {
						s.store[set.Key] = set.Value
					}
				}
			}
		}
	}
}

func (s *KeyValue) set() {
	for set := range s.pending {
		s.storeLock.Lock()
		if set.Value == "" {
			delete(s.store, set.Key)
		} else {
			s.store[set.Key] = set.Value
		}
		s.storeLock.Unlock()

		s.pendingPersist <- set
	}
}

func (s *KeyValue) persistDelta() {
	ticker := time.NewTicker(time.Second)
	for t := range ticker.C {
		func(s *KeyValue) {
			length := len(s.pendingPersist)
			if length == 0 {
				return
			}

			buffer := make([]*set, length)
			for i := 0; i < length; i++ {
				buffer[i] = <-s.pendingPersist
			}

			deltaPath := path.Join(s.LogDir, fmt.Sprintf("%d-delta", t.UnixNano()))
			f, err := os.Create(deltaPath)
			if err != nil {
				log.Printf("Could not create file %s, failed with error: %v\n", deltaPath, err)
				return
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			defer w.Flush()

			data, err := json.Marshal(buffer)
			if err != nil {
				log.Printf("Could not marshall delta log, with error: %v\n", err)
			}
			w.Write(data)
			if err != nil {
				log.Printf("Could not write data failed, with error: %v\n", err)
			}
		}(s)
	}
}

func (s *KeyValue) persistBase() {
	ticker := time.NewTicker(time.Minute)
	for t := range ticker.C {
		func(s *KeyValue) {
			basePath := path.Join(s.LogDir, fmt.Sprintf("%d-base", t.UnixNano()))
			f, err := os.Create(basePath)
			if err != nil {
				log.Printf("Could not create file %s, failed with error: %v\n", basePath, err)
				return
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			defer w.Flush()

			s.storeLock.RLock()
			data, err := json.Marshal(s.store)
			s.storeLock.RUnlock()
			if err != nil {
				log.Printf("Could not marshall delta log, with error: %v\n", err)
			}
			w.Write(data)
			if err != nil {
				log.Printf("Could not write data failed, with error: %v\n", err)
			}
			go s.deleteOldPersistence(t.UnixNano())
		}(s)
	}
}

func (s *KeyValue) deleteOldPersistence(epoch int64) {
	entries, err := ioutil.ReadDir(s.LogDir)
	if err != nil {
		log.Printf("Error reading log directory: %v", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.LastIndex(name, "-base") >= 0 || strings.LastIndex(name, "-delta") >= 0 {
			split := strings.Split(name, "-")
			if len(split) == 2 {
				touch, err := strconv.ParseInt(split[0], 10, 64)
				if err == nil && touch < epoch {
					os.Remove(path.Join(s.LogDir, name))
				}
			}
		}
	}
}

func (s *KeyValue) Get(key string) (int, string) {
	if s.store == nil {
		log.Printf("KeyValue Store is not initialized\n")
		return -1, ""
	}

	s.storeLock.RLock()
	value, present := s.store[key]
	s.storeLock.RUnlock()
	if present {
		return 0, value
	}
	return 1, ""
}

func (s *KeyValue) Set(key string, value string) (int, string) {
	status, oldValue := s.Get(key)

	s.pending <- &set{Key: key, Value: value}

	return status, oldValue
}
