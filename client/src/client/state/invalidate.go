package state

import (
	"log"

	"client/keyvalue"
	"client/thrift/pool"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Invalidate struct {
	Shards   []keyvalue.Shard
	File     *keyvalue.File
	CallBack func(*Invalidate, *StateMachine)
}

func (i Invalidate) Run(sm *StateMachine) {
	clientPool := pool.NewClientPool(thrift.NewTBufferedTransportFactory(10000), thrift.NewTBinaryProtocolFactoryDefault())
	defer clientPool.CloseAll()

	for _, shard := range i.Shards {

		client, err := clientPool.GetClient(shard.IP)
		if err != nil {
			log.Println("Error opening connection to ", shard.IP, err)
			continue
		}

		err = client.Remove(shard.Id)
		if err != nil {
			log.Println("Error during invalidate", shard.Id, err)
		}
	}

	if i.CallBack != nil {
		i.CallBack(&i, sm)
	}
}
