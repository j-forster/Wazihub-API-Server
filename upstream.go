package main

import (
	"log"
	"strings"
	"time"

	wazihub "github.com/j-forster/Wazihub-API"
	"github.com/j-forster/Wazihub-API/mqtt"
)

var upstream *mqtt.Queue
var upstreamId = wazihub.CurrentDeviceId()

var retries = []time.Duration{
	5 * time.Second,
	10 * time.Second,
	20 * time.Second,
	60 * time.Second,
}

func Upstream(addr string) {

	if !strings.ContainsRune(addr, ':') {
		addr = addr + ":1883"
	}

	upstream = mqtt.NewQueue("upstream")
	nretry := 0

	for {
		log.Printf("[UP   ] Dialing Upstream at %q...\n", addr)
		client, err := mqtt.Dial(addr, upstreamId, false, nil, nil)
		if err != nil {
			log.Printf("[UP   ] Error: %v\n", err)
			duration := retries[nretry]
			log.Printf("[UP   ] Waiting %s before retry.\n", duration)
			time.Sleep(duration)
			nretry++
			if nretry == len(retries) {
				nretry = len(retries) - 1
			}
			continue
		}
		client.Subscribe("devices/"+upstreamId+"/actuators/#", 0)
		log.Printf("[UP   ] Connected.\n")
		upstream.ServeWriter(client)

		for msg := range client.Message() {
			log.Printf("[UP   ] Recieved \"%s\" QoS:%d len:%d\n", msg.Topic, msg.QoS, len(msg.Data))
			mqttServer.Publish(client, msg)
		}

		log.Printf("[UP   ] Disconnected: %v\n", client.Error)
	}
}
