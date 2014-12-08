package state

import (
	"crypto/aes"
	"crypto/cipher"
	"sync/atomic"
	"time"

	"client/settings"

	"github.com/somethingnew2-0/go-erasure"
)

type StateMachine struct {
	Options     *settings.Options
	ErasureCode *erasure.Code
	Cipher      cipher.Block
	Recovered   bool // If user already existed, recover files
	Initialized bool // If user didn't exist, upload initial files before watching
	states      chan State
	workers     uint32
}

type State interface {
	Run(sm *StateMachine)
}

func NewStateMachine(opts *settings.Options) *StateMachine {
	cipher, _ := aes.NewCipher(opts.Hash)
	return &StateMachine{
		Options:     opts,
		ErasureCode: erasure.NewCode(settings.M, settings.K, settings.BlockSize),
		Cipher:      cipher,
		states:      make(chan State, settings.MaxStates),
		workers:     0,
	}
}

func (sm StateMachine) Run() {
	for i := 0; i < settings.MinimumWorkers; i++ {
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

func (sm StateMachine) Add(s State) {
	sm.states <- s
}

func (sm StateMachine) SpawnWorker() {
	go func() {
		atomic.AddUint32(&sm.workers, 1)
		for state := range sm.states {
			state.Run(&sm)

			// Retire any workers over the threshold, once free
			workers := atomic.LoadUint32(&sm.workers)
			if settings.MinimumWorkers < int(workers) {
				if atomic.CompareAndSwapUint32(&sm.workers, workers, workers-1) {
					return
				}
			}
		}
	}()
}
