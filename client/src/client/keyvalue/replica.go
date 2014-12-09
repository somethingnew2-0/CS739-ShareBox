package keyvalue

import (
	"encoding/json"
	"errors"
)

type Replica struct {
}

func InitReplicaKV() *KeyValue {
	return Init("log/replica")
}

func (kv KeyValue) GetReplica(fileId string) (*Replica, error) {
	status, replicaJson := kv.Get(fileId)
	if status != 0 {
		return nil, errors.New("Replica doesn't exist in the key value store")
	}
	replica := &Replica{}
	json.Unmarshal([]byte(replicaJson), &replica)
	return replica, nil
}

func (kv KeyValue) SetReplica(fileId string, replica *Replica) error {
	replicaJson, err := json.Marshal(replica)
	if err != nil {
		return err
	}
	status, _ := kv.Set(fileId, string(replicaJson))

	if status != 0 {
		return errors.New("Error in setting Replica in the key value store")
	}
	return nil
}
