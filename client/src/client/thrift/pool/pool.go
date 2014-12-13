package pool

import (
	"crypto/tls"
	"fmt"

	"client/settings"
	"client/thrift/replica"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type ClientPool struct {
	transports       []*thrift.TTransport
	clients          map[string]*replica.ReplicatorClient
	transportFactory thrift.TTransportFactory
	protocolFactory  thrift.TProtocolFactory
}

func NewClientPool(tf thrift.TTransportFactory, pf thrift.TProtocolFactory) *ClientPool {
	return &ClientPool{
		transports:       make([]*thrift.TTransport, 0),
		clients:          make(map[string]*replica.ReplicatorClient),
		transportFactory: tf,
		protocolFactory:  pf,
	}
}
func (tp ClientPool) GetClient(ip string) (*replica.ReplicatorClient, error) {
	client := tp.clients[ip]
	if client == nil {
		addr := fmt.Sprintf("%s:%d", ip, settings.ClientPort)
		var err error
		var socket thrift.TTransport
		if settings.ClientTLS {
			cfg := new(tls.Config)
			cfg.InsecureSkipVerify = true
			socket, err = thrift.NewTSSLSocket(addr, cfg)
		} else {
			socket, err = thrift.NewTSocket(addr)
		}
		if err != nil {
			return nil, err
		}
		transport := tp.transportFactory.GetTransport(socket)
		if err := transport.Open(); err != nil {
			return nil, err
		}
		// TODO: Actually figure out caching later
		tp.transports = append(tp.transports, &transport)
		client = replica.NewReplicatorClientFactory(transport, tp.protocolFactory)
		tp.clients[ip] = client
	}
	client.Ping()
	return client, nil
}

func (tp ClientPool) CloseAll() {
	for _, transport := range tp.transports {
		(*transport).Close()
	}
}
