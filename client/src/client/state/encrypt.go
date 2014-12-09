package state

import (
	"client/keyvalue"
)

type Encrypt struct {
	File      *keyvalue.File
	Plaintext []byte
}

func (e Encrypt) Run(sm *StateMachine) {
	encode := &Encode{File: e.File, Ciphertext: make([]byte, len(e.Plaintext))}
	if len(e.Plaintext) > 0 {
		sm.Cipher.Encrypt(encode.Ciphertext, e.Plaintext)
	}
	sm.Add(encode)
}
