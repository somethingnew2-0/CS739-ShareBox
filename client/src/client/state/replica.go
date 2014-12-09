package state

import (
	"crypto/tls"
	"fmt"
	"log"

	"client/settings"
	"client/thrift/replica"

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
	fmt.Printf("%T\n", transport)
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

	return nil
}

func (rh ReplicaHandler) Modify(r *replica.Replica) error {
	return nil
}

func (rh ReplicaHandler) Remove(shardId string) error {
	return nil
}

func (rh ReplicaHandler) Download(shardId string) (*replica.Replica, error) {
	return nil, nil
}
