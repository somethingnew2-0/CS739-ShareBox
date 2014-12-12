package state

import (
	"client/keyvalue"
)

type Encrypt struct {
	Create    bool // Is this a create or modify
	File      *keyvalue.File
	Plaintext []byte
}

func (e Encrypt) Run(sm *StateMachine) {
	encode := &Encode{
		Create:     e.Create,
		File:       e.File,
		Ciphertext: make([]byte, len(e.Plaintext)),
	}
	if len(e.Plaintext) > 0 {
		sm.Cipher.Encrypt(encode.Ciphertext, e.Plaintext)
	}
	sm.Add(encode)
}
