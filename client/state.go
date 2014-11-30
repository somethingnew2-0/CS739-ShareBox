package main

import (
	"sync/atomic"
	"time"
)

const minimumWorkers int = 16
const maxStates int = 256

type StateMachine struct {
	ClientId string
	states   chan State
	workers  uint32
}

type State interface {
	Run(sm *StateMachine)
}

func NewStateMachine(clientId string) *StateMachine {
	return &StateMachine{
		ClientId: clientId,
		states:   make(chan State, maxStates),
		workers:  0,
	}
}

func (sm *StateMachine) Run() {
	for i := 0; i < minimumWorkers; i++ {
		sm.SpawnWorker()
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	for _ = range ticker.C {
		// Check for deadlocks and spawn more workers when needed
		if len(sm.states) == cap(sm.states) {
			sm.SpawnWorker()
		}
	}
}

func (sm *StateMachine) SpawnWorker() {
	go func() {
		atomic.AddUint32(&sm.workers, 1)
		for state := range sm.states {
			state.Run(sm)

			// Retire any workers over the threshold, once free
			workers := atomic.LoadUint32(&sm.workers)
			if minimumWorkers < int(workers) {
				if atomic.CompareAndSwapUint32(&sm.workers, workers, workers-1) {
					return
				}
			}
		}
	}()
}
