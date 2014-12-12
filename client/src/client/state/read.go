package state

import (
	"io/ioutil"
	"log"
	"math"
	"os"

	"client/keyvalue"
	"client/settings"
)

type Read struct {
	Create bool // Is this a create or modify
	Path   string
	Info   os.FileInfo
}

func (r Read) Run(sm *StateMachine) {
	blocks := math.Ceil(float64(r.Info.Size()) / float64(settings.BlockSize))

	encrypt := &Encrypt{
		Create: r.Create,
		File: &keyvalue.File{
			Name: r.Path,
			// Use encoded file size
			EncodedSize:   int64(blocks) * settings.BlockSize,
			UnencodedSize: r.Info.Size(),
		},
	}

	if f, err := os.Open(r.Path); err == nil {
		zeroBytes := (int64(blocks) * settings.BlockSize) - r.Info.Size()
		if data, err := ioutil.ReadAll(f); err == nil {
			encrypt.Plaintext = append(data, make([]byte, zeroBytes)...)
		} else {
			log.Println("Cannot read file: ", r.Path, err)
		}
		sm.Add(encrypt)
	} else {
		log.Println("Cannot open file: ", r.Path, err)
	}
}
