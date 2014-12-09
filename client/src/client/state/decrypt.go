package state

import (
	"crypto/sha256"
	"crypto/subtle"
	"io/ioutil"
	"log"

	"client/keyvalue"
)

type Decrypt struct {
	File       *keyvalue.File
	Ciphertext []byte
}

func (d Decrypt) Run(sm *StateMachine) {
	fileHash := sha256.New()
	fileHash.Write(d.Ciphertext)
	if subtle.ConstantTimeCompare(fileHash.Sum(nil), []byte(d.File.Hash)) == 1 {
		plaintext := make([]byte, len(d.Ciphertext))
		sm.Cipher.Decrypt(plaintext, d.Ciphertext)
		err := sm.Files.SetFile(d.File.Name, d.File)
		if err != nil {
			log.Println("Error saving metadata for recovered file: ", d.File.Name, err)
		}

		err = ioutil.WriteFile(d.File.Name, plaintext[:d.File.UnencodedSize], 0666)
		if err != nil {
			log.Println("Error saving recovered file: ", d.File.Name, err)
		}
	} else {
		log.Println("Recovered file was corrupted! ", d.File.Name)
	}
}
