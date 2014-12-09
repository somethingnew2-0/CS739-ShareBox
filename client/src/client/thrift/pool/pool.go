package pool

import (
	"crypto/tls"
	"fmt"

	"client/settings"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type TransportPool struct {
	transports       map[string]thrift.TTransport
	transportFactory thrift.TTransportFactory
}

func NewTransportPool(tf thrift.TTransportFactory) *TransportPool {
	return &TransportPool{
		transports:       make(map[string]thrift.TTransport),
		transportFactory: tf,
	}
}
func (tp TransportPool) GetTransport(ip string) (thrift.TTransport, error) {
	transport := tp.transports[ip]
	if transport == nil {
		addr := fmt.Sprintf("%s:%d", ip, settings.ClientPort)
		var err error
		if settings.ClientTLS {
			cfg := new(tls.Config)
			cfg.InsecureSkipVerify = true
			transport, err = thrift.NewTSSLSocket(addr, cfg)
		} else {
			transport, err = thrift.NewTSocket(addr)
		}
		if err != nil {
			return nil, err
		}
		transport = tp.transportFactory.GetTransport(transport)
		if err := transport.Open(); err != nil {
			return nil, err
		}
		tp.transports[ip] = transport
	}
	return transport, nil
}

func (tp TransportPool) CloseAll() {
	for _, transport := range tp.transports {
		transport.Close()
	}
}
