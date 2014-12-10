package keyvalue

import (
	"encoding/json"
	"errors"
)

type Shard struct {
	Id     string `json:"id"`
	Hash   string `json:"hash"`
	Offset int64  `json:"offset"`
	IP     string `json:"ip"`
	Size   int64  `json:"size"`
}

type Block struct {
	Id          string  `json:"id"`
	Hash        string  `json:"hash"`
	BlockOffset int64   `json:"blockOffset"`
	Shards      []Shard `json:"shards"`
}

type File struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Hash          string  `json:"hash"`
	EncodedSize   int64   `json:"size"`
	UnencodedSize int64   `json:"originalSize"`
	Blocks        []Block `json:"blocks"`
}

func InitFileKV() *KeyValue {
	return Init("log/files")
}

func (kv KeyValue) GetFile(path string) (*File, error) {
	status, fileJson := kv.Get(path)
	if status != 0 || fileJson == "" {
		return nil, errors.New("File doesn't exist in the key value store")
	}
	file := &File{}
	json.Unmarshal([]byte(fileJson), &file)
	return file, nil
}

func (kv KeyValue) SetFile(path string, file *File) error {
	fileJson := []byte("")
	if file != nil {
		var err error
		fileJson, err = json.Marshal(file)
		if err != nil {
			return err
		}
	}
	kv.Set(path, string(fileJson))
	return nil
}
