package state

import (
	"crypto/sha256"

	"client/keyvalue"
	"client/settings"
)

type Checksum struct {
	Create        bool
	File          *keyvalue.File
	EncodedBlocks [][]byte
}

func (c Checksum) Run(sm *StateMachine) {
	file := c.File
	fileHash := sha256.New()
	for i, block := range c.EncodedBlocks {
		fileHash.Write(block)

		blockHash := sha256.New()
		blockHash.Write(block)

		file.Blocks = append(file.Blocks, keyvalue.Block{
			Hash:        blockHash.Sum(nil),
			BlockOffset: int64(settings.BlockSize * i)})

		for s := 0; s < settings.M; s++ {
			shard := block[s*settings.ShardLength : (s+1)*settings.ShardLength]
			shardHash := sha256.New()
			shardHash.Write(shard)
			file.Blocks[i].Shards = append(file.Blocks[i].Shards, keyvalue.Shard{
				Hash:   shardHash.Sum(nil),
				Offset: int64(s),
				Size:   int64(settings.ShardLength),
			})
		}
	}

	file.Hash = fileHash.Sum(nil)

	if c.Create {
		sm.Add(&Create{EncodedBlocks: c.EncodedBlocks, File: file})
	} else {
		sm.Add(&Update{EncodedBlocks: c.EncodedBlocks, File: file})
	}
}
