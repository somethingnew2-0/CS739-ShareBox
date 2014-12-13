package state

import (
	"crypto/subtle"
	"fmt"
	"log"
	"time"

	"client/keyvalue"
	"client/settings"
	"client/util"
)

type Update struct {
	EncodedBlocks [][]byte
	File          *keyvalue.File
}

func (u Update) Run(sm *StateMachine) {
	f, err := sm.Files.GetFile(u.File.Name)

	// Updates usually follow a Create.  Wait a bit for it to complete.
	// If the Create never happens or takes too long to finish, we'll do it ourselves.
	failures := 0
	ticker := time.NewTicker(time.Second)
	for err != nil {
		failures++
		log.Println("File not created yet trying again ", failures, u.File.Name)
		// If we reach the timeout, just create the file instead
		if failures > settings.UpdateTimeout {
			sm.Add(&Create{EncodedBlocks: u.EncodedBlocks, File: u.File})
			return
		}
		<-ticker.C
		f, err = sm.Files.GetFile(u.File.Name)
	}

	if subtle.ConstantTimeCompare([]byte(u.File.Hash), []byte(f.Hash)) == 1 {
		log.Println("No changes were actaully detected with file update")
		return
	}

	updated := &keyvalue.File{
		Id:            f.Id,
		Name:          u.File.Name,
		Hash:          u.File.Hash,
		EncodedSize:   u.File.EncodedSize,
		UnencodedSize: u.File.UnencodedSize,
		Blocks:        make([]keyvalue.Block, 0),
	}

	invalidate := make([]keyvalue.Shard, 0)

	for i, block := range u.File.Blocks {
		if len(f.Blocks) < i {
			if subtle.ConstantTimeCompare([]byte(block.Hash), []byte(f.Blocks[i].Hash)) == 1 {
				// The block hasn't changed
				continue
			}
			// The block has changed
			block.Id = f.Blocks[i].Id
			invalidate = append(invalidate, f.Blocks[i].Shards...)
		}
		updated.Blocks = append(updated.Blocks, block)
	}

	// Check if the file shrunk
	if len(u.File.Blocks) < len(f.Blocks) {
		// Call block delete for discarded blocks?
		// Does the server automatically discard blocks based on file size?
		// How do the replica client know this?

		// Invalidate any blocks larger than the
		for _, block := range f.Blocks[len(u.File.Blocks):] {
			invalidate = append(invalidate, block.Shards...)
		}
	}

	resp, err := util.Post(sm.Options, fmt.Sprintf("client/%s/file/update", sm.Options.ClientId), updated)
	if err != nil {
		log.Println("Error adding file", err)
		return
	}

	if resp["allowed"].(bool) {
		file := u.File
		file.Id = resp["id"].(string)
		blockIds := resp["blocks"].([]interface{})
		for i, blockId := range blockIds {
			file.Blocks[i].Id = blockId.(string)
		}

		clients := resp["clients"].([]interface{})
		for _, c := range clients {
			client := c.(map[string]interface{})
			blockId := client["blockId"]
			offset := int(client["offset"].(float64))
			for _, block := range file.Blocks {
				if block.Id == blockId {
					block.Shards[offset].Id = client["id"].(string)
					// TODO Validate this an IP address using net.IP
					block.Shards[offset].IP = client["IP"].(string)
					break
				}
			}
		}
		sm.Add(&Upload{
			EncodedBlocks: u.EncodedBlocks,
			File:          file,
		})
		sm.Add(&Invalidate{
			Shards: invalidate,
			// Not actually used
			File: u.File,
		})
	} else {
		log.Println("File upload not allowed", u.File.Name)
		return
	}
}
