package main

import (
	"log"

	"github.com/j-forster/Wazihub-API/mqtt"
)

var upstream chan *mqtt.Message

func Upstream(addr string) {

	upstream = make(chan *mqtt.Message)

	log.Printf("[UP   ] Dialing Upstream at %q...\n", addr)
	conn, err := mqtt.Dial(addr, "Mario", true, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		msg := <-upstream
		conn.Publish(nil, msg)
	}
}
