package state

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	"client/settings"
	"client/thrift/replica"
	"client/util"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Remove struct {
	Path string
}

func (r Remove) Run(sm *StateMachine) {
	file, err := sm.Files.GetFile(r.Path)
	var resp map[string]interface{}
	if err == nil {
		resp, err = util.Post(fmt.Sprintf("client/%s/file/remove", sm.Options.ClientId), url.Values{"name": {r.Path}})
		if err != nil {
			log.Println("Error removing file", err)
			return
		}
	} else {
		resp, err = util.Post(fmt.Sprintf("client/%s/file/remove", sm.Options.ClientId), url.Values{"id": {file.Id}, "name": {r.Path}, "size": {string(file.EncodedSize)}})
		if err != nil {
			log.Println("Error removing file", err)
			return
		}
	}

	if resp["allowed"].(bool) {
		shards := resp["clients"].([]map[string]string)
		for _, shard := range shards {
			transports := make(map[string]thrift.TTransport)
			transportFactory := thrift.NewTBufferedTransportFactory(1024)
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
			// TODO Validate this an IP address using net.IP
			transport := transports[shard["IP"]]
			if transport == nil {
				addr := fmt.Sprintf("%s:%d", shard["IP"], settings.ClientPort)
				if settings.ClientTLS {
					cfg := new(tls.Config)
					cfg.InsecureSkipVerify = true
					transport, err = thrift.NewTSSLSocket(addr, cfg)
				} else {
					transport, err = thrift.NewTSocket(addr)
				}
				if err != nil {
					log.Println("Error opening socket:", err)
					break
				}
				transport = transportFactory.GetTransport(transport)
				defer transport.Close()
				if err := transport.Open(); err != nil {
					break
				}
				transports[shard["IP"]] = transport
			}
			client := replica.NewReplicatorClientFactory(transport, protocolFactory)
			client.Ping()

			iv, err := client.Remove(shard["id"])

			if iv != nil {
				log.Println("Invalid operation:", iv)
			} else if err != nil {
				log.Println("Error during remove", err)
			}
		}

		if file != nil {
			resp, err = util.Post(fmt.Sprintf("file/%s/delete", file.Id), url.Values{"clientId": {sm.Options.ClientId}})
			if err != nil {
				log.Println("Error deleting file", err)
				return
			}
			if !resp["success"].(bool) {
				log.Println("Deleting file from server was unsucessful")
			}
		}
	}
}
