package state

import (
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"client/keyvalue"
	"client/settings"
	"client/thrift/replica"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Replica struct{}

func (r Replica) Run(sm *StateMachine) {
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(10000)
	var transport thrift.TServerTransport
	var err error
	addr := fmt.Sprintf(":%d", settings.ClientPort)
	if settings.ClientTLS {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			log.Println("Could not start TLS replica server on client", err)
			return
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	if err != nil {
		log.Println("Could not start replica server on client", err)
		return
	}
	handler := &ReplicaHandler{StateMachine: sm}
	processor := replica.NewReplicatorProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	server.Serve()
}

type ReplicaHandler struct {
	StateMachine *StateMachine
}

func (rh ReplicaHandler) Ping() error {
	return nil
}

func (rh ReplicaHandler) Add(r *replica.Replica) error {
	resp, err := util.Post(rh.StateMachine.Options, fmt.Sprintf("shard/%s/validate", r.ShardId), map[string]string{"receiverId": rh.StateMachine.Options.ClientId, "ownerId": r.ClientId})
	if err != nil {
		return err
	}
	if resp["error"] != nil {
		return errors.New(fmt.Sprintf("Error validating shard %s %s", resp["error"], resp["message"]))
	}
	if !resp["accept"].(bool) {
		return errors.New("Adding this shard to the replica is not allowed")
	}

	shardHash := sha256.New()
	shardHash.Write(r.Shard)
	// if subtle.ConstantTimeCompare(shardHash.Sum(nil), []byte(resp["hash"].(string))) == 0 ||
	if subtle.ConstantTimeCompare(shardHash.Sum(nil), []byte(r.ShardHash)) == 0 {
		return errors.New("Shard didn't match shard hash")
	}
	replica := &keyvalue.Replica{
		ShardHash:   r.ShardHash,
		ShardOffset: r.ShardOffset,
		ShardId:     r.ShardId,
		BlockId:     r.BlockId,
		FileId:      r.FileId,
		ClientId:    r.ClientId,
	}
	err = os.MkdirAll(path.Dir(getPath(replica.ShardId)), 0777)
	if err != nil {
		log.Println("Error create replica directory ", err)
		return err
	}
	err = ioutil.WriteFile(getPath(replica.ShardId), r.Shard, 0666)
	if err != nil {
		log.Println("Error writing shard file ", err)
		return err
	}
	err = rh.StateMachine.Replicas.SetReplica(r.ShardId, replica)
	if err != nil {
		log.Println("Error setting shard in replica metadata ", err)
		return err
	}

	return nil
}

func (rh ReplicaHandler) Remove(shardId string) error {
	r, err := rh.StateMachine.Replicas.GetReplica(shardId)
	if err != nil {
		return err
	}

	resp, err := util.Post(rh.StateMachine.Options, fmt.Sprintf("shard/%s/invalidate", r.ShardId), map[string]string{"receiverId": rh.StateMachine.Options.ClientId, "ownerId": r.ClientId})
	if err != nil {
		return err
	}
	if resp["error"] != nil {
		return errors.New(fmt.Sprintf("Error invalidating shard %s %s", resp["error"], resp["message"]))
	}
	if !resp["delete"].(bool) {
		return errors.New("Deleting this shard to the replica is not allowed")
	}

	err = rh.StateMachine.Replicas.SetReplica(shardId, nil)
	removeErr := os.Remove(getPath(shardId))
	if err != nil {
		return err
	}
	if removeErr != nil {
		return removeErr
	}
	return nil
}

func (rh ReplicaHandler) Download(shardId string) (*replica.Replica, error) {
	r, err := rh.StateMachine.Replicas.GetReplica(shardId)
	if err != nil {
		return nil, err
	}
	replica := &replica.Replica{
		ShardHash:   r.ShardHash,
		ShardOffset: r.ShardOffset,
		ShardId:     r.ShardId,
		BlockId:     r.BlockId,
		FileId:      r.FileId,
		ClientId:    r.ClientId,
	}
	if _, err := os.Stat(getPath(shardId)); err == nil {
		if f, err := os.Open(getPath(shardId)); err == nil {
			shard, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			shardHash := sha256.New()
			shardHash.Write(shard)
			if subtle.ConstantTimeCompare(shardHash.Sum(nil), []byte(replica.ShardHash)) == 0 {
				return nil, errors.New("Shard didn't match shard hash")
			}

			replica.Shard = shard
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	return replica, nil
}

func getPath(shardId string) string {
	firstTwo := shardId[:2]
	secondTwo := shardId[2:4]
	fileName := shardId[4:]

	return fmt.Sprintf("%s/%s/%s/%s", settings.ReplicasPath, firstTwo, secondTwo, fileName)
}
