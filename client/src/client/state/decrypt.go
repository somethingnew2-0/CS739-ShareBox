package state

import (
	"client/keyvalue"
)

type Decrypt struct {
	File          *keyvalue.File
	DecodedBlocks [][]byte
}

func (d Decrypt) Run(sm *StateMachine) {
}
