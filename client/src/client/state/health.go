package state

import (
	"fmt"
	"log"
	"time"

	"client/settings"

	"github.com/armon/consul-api"
)

type Health struct{}

func (h Health) Run(sm *StateMachine) {
	config := consulapi.DefaultConfig()
	config.Address = settings.ConsulAddress

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Println("Could not connect to consul: ", err)
	}
	agent := client.Agent()
	info, err := agent.Self()
	if err != nil {
		log.Println("Failed to identify ourselved to consul: ", err)
	}

	// Join ourself
	addr := info["Config"]["AdvertiseAddr"].(string)
	err = agent.Join(addr, false)
	if err != nil {
		log.Println("Failed to join consul", err)
	}

	defer agent.ForceLeave(addr)

	registration := &consulapi.AgentServiceRegistration{}
	registration.ID = sm.Options.ClientId
	registration.Name = "client"
	registration.Port = settings.ClientPort
	registration.Tags = []string{"client"}
	registration.Check = &consulapi.AgentServiceCheck{
		TTL:      "10m",
		Interval: "1m",
	}

	agent.ServiceRegister(registration)
	defer agent.FailTTL(fmt.Sprintf("service:%s", registration.ID), "Shutdown")
	agent.PassTTL(fmt.Sprintf("service:%s", registration.ID), "Still alive!")

	ticker := time.NewTicker(time.Minute)
	for _ = range ticker.C {
		agent.PassTTL(fmt.Sprintf("service:%s", registration.ID), "Still alive!")
	}
}
