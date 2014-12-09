package state

import (
	"io/ioutil"
	"log"

	"client/keyvalue"
)

type Decrypt struct {
	File       *keyvalue.File
	Ciphertext []byte
}

func (d Decrypt) Run(sm *StateMachine) {
	plaintext := make([]byte, len(d.Ciphertext))
	sm.Cipher.Decrypt(plaintext, d.Ciphertext)

	err := sm.Files.SetFile(d.File.Name, d.File)
	if err != nil {
		log.Println("Error saving metadata for recovered file: ", err)
	}
	err = ioutil.WriteFile(d.File.Name, plaintext, 0666)
	if err != nil {
		log.Println("Error saving recovered file: ", err)
	}
}
