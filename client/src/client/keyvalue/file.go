package keyvalue

import (
	"encoding/json"
	"errors"
	"net"
)

type Shard struct {
	Id     string `json:"id"`
	Hash   string `json:"hash"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
	IP     net.IP `json:"ip"`
}

type Block struct {
	Id          string  `json:"id"`
	Hash        string  `json:"hash"`
	BlockOffset int     `json:"blockOffset"`
	Shards      []Shard `json:"shards"`
}

type File struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Hash   string  `json:"hash"`
	Size   int64   `json:"size"`
	Blocks []Block `json:"blocks"`
}

func InitFileKV() *KeyValue {
	return Init("log/files")
}

func (kv KeyValue) GetFile(path string) (*File, error) {
	status, fileJson := kv.Get(path)
	if status != 0 {
		return nil, errors.New("File doesn't exist in the key value store")
	}
	file := &File{}
	json.Unmarshal([]byte(fileJson), &file)
	return file, nil
}

func (kv KeyValue) SetFile(path string, file *File) error {
	fileJson, err := json.Marshal(file)
	if err != nil {
		return err
	}
	status, _ := kv.Set(path, string(fileJson))

	if status != 0 {
		return errors.New("Error in setting File in the key value store")
	}
	return nil
}
