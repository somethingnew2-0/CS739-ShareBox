package state

import (
	"client/keyvalue"
)

type Encrypt struct {
	File      *keyvalue.File
	Plaintext []byte
}

func (e Encrypt) Run(sm *StateMachine) {
	encode := &Encode{Ciphertext: make([]byte, len(e.Plaintext))}
	sm.Cipher.Encrypt(encode.Ciphertext, e.Plaintext)
	sm.Add(encode)
}
