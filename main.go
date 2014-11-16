package main

import (
	"bytes"
	"log"
	"math/rand"

	"github.com/somethingnew2-0/go-erasure"
)

func corrupt(source, errList []byte, shardLength int) []byte {
	corrupted := make([]byte, len(source))
	copy(corrupted, source)
	for _, err := range errList {
		for i := 0; i < shardLength; i++ {
			corrupted[int(err)*shardLength+i] = 0x00
		}
	}
	return corrupted
}

func main() {
	m := 12
	k := 8
	shardLength := 16       // Length of a shard
	size := k * shardLength // Length of the data blob to encode

	code := erasure.NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 2, 3, 4}

	corrupted := corrupt(append(source, encoded...), errList, shardLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		log.Fatal("Source was not sucessfully recovered with 4 errors")
	}
}
