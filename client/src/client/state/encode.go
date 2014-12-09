package state

import (
	"client/keyvalue"
	"client/settings"
)

type Encode struct {
	File       *keyvalue.File
	Ciphertext []byte
}

func (e Encode) Run(sm *StateMachine) {
	blocks := len(e.Ciphertext) / settings.BlockSize
	checksum := &Checksum{File: e.File, EncodedBlocks: make([][]byte, blocks)}
	if len(e.Ciphertext) > 0 {
		for i, _ := range checksum.EncodedBlocks {
			encryptedBlock := e.Ciphertext[i*settings.BlockSize : (i+1)*settings.BlockSize]
			checksum.EncodedBlocks[i] = append(encryptedBlock, sm.ErasureCode.Encode(encryptedBlock)...)
		}
	}
	sm.Add(checksum)
}
