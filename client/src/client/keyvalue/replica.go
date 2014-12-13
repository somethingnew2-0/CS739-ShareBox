package keyvalue

import (
	"encoding/json"
	"errors"
)

type Replica struct {
	ShardHash   []byte `json:"shardHash"`
	ShardOffset int32  `json:"shardOffset"`
	ShardId     string `json:"shardId"`
	BlockId     string `json:"blockId"`
	FileId      string `json:"fileId"`
	ClientId    string `json:"clientId"`
}

func InitReplicaKV() *KeyValue {
	return Init("log/replica")
}

func (kv KeyValue) GetReplica(shardId string) (*Replica, error) {
	status, replicaJson := kv.Get(shardId)
	if status != 0 || replicaJson == "" {
		return nil, errors.New("Replica doesn't exist in the key value store")
	}
	replica := &Replica{}
	json.Unmarshal([]byte(replicaJson), &replica)
	return replica, nil
}

func (kv KeyValue) SetReplica(shardId string, replica *Replica) error {
	replicaJson := []byte("")
	if replica != nil {
		var err error
		replicaJson, err = json.Marshal(replica)
		if err != nil {
			return err
		}
	}
	kv.Set(shardId, string(replicaJson))
	return nil
}
