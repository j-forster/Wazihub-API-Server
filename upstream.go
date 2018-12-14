package main

import (
	"log"
	"time"

	"github.com/j-forster/Wazihub-API/mqtt"
)

var upstream *mqtt.Queue

func Upstream(addr string) {

	upstream = mqtt.NewQueue("upstream")

	for {
		log.Printf("[UP   ] Dialing Upstream at %q...\n", addr)
		client, err := mqtt.Dial(addr, "upstream", false, nil, nil)
		if err != nil {
			log.Printf("[UP   ] Error: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[UP   ] Connected.\n")
		upstream.ServeWriter(client)

		for msg := range client.Message() {
			mqttServer.Publish(client, msg)
		}

		log.Printf("[UP   ] Disconnected: %v\n", client.Error)
	}
}
