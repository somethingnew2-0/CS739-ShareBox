package state

import (
	"log"

	"client/keyvalue"
	"client/settings"
)

type Decode struct {
	File          *keyvalue.File
	EncodedBlocks [][]byte
	BlockErrs     [][]byte
}

func (d Decode) Run(sm *StateMachine) {
	decrypt := &Decrypt{
		File: d.File,
	}
	for b, block := range d.EncodedBlocks {
		if len(d.BlockErrs[b]) > settings.M-settings.K {
			log.Println("Too many errors cannot recover file: ", d.File.Name)
			return
		}
		decrypt.Ciphertext = append(decrypt.Ciphertext, sm.ErasureCode.Decode(block, d.BlockErrs[b])...)
	}
	sm.Add(decrypt)
}
