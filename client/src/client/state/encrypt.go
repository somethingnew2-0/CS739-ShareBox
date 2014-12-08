package state

import (
	"crypto/aes"
)

type Encrypt struct {
	Plaintext []byte
}

func (e Encrypt) Run(sm *StateMachine) {
	encode := &Encode{Ciphertext: make([]byte, len(e.Plaintext))}
	c, _ := aes.NewCipher(sm.Options.Hash)
	c.Encrypt(encode.Ciphertext, e.Plaintext)
}
