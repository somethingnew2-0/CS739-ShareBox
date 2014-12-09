package state

import (
	"io/ioutil"
	"math"
	"os"

	"client/keyvalue"
	"client/settings"
)

type Create struct {
	Path string
	Info os.FileInfo
}

func (c Create) Run(sm *StateMachine) {
	encrypt := &Encrypt{}

	blocks := math.Ceil(float64(c.Info.Size()) / float64(settings.BlockSize))

	// Use encoded file size
	encrypt.File = &keyvalue.File{Name: c.Path, EncodedSize: int64(blocks) * settings.BlockSize, UnencodedSize: c.Info.Size()}
	if f, err := os.Open(c.Path); err == nil {
		zeroBytes := (int64(blocks) * settings.BlockSize) - c.Info.Size()
		data, _ := ioutil.ReadAll(f)
		encrypt.Plaintext = append(data, make([]byte, zeroBytes)...)
	}
	sm.Add(encrypt)
}
