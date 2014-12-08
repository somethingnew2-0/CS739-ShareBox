package state

import (
	"crypto/sha256"
	"log"
)

type Checksum struct {
	EncodedCiphertext []byte
}

func (c Checksum) Run(sm *StateMachine) {
	hash := sha256.New()

	hash.Write(c.EncodedCiphertext)

	log.Printf("Checksum is %x\n", hash.Sum(nil))
}
