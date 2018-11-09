package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/j-forster/Wazihub-API/mqtt"
	"github.com/j-forster/Wazihub-API/tools"
)

var mqttHandler = &MQTTHandler{}
var mqttServer = mqtt.NewServer(nil, mqttHandler)

func init() {
	go mqttServer.Run()
}

func ListenAndServerMQTT() {

	log.Println("[MQTT ] MQTT Server at \":1883\".")
	listener, err := net.Listen("tcp", ":1883")
	if err != nil {
		log.Fatalln("[MQTT ] Error:\n", err)
	}

	for {

		conn, err := listener.Accept()
		if err == nil {

			go mqttServer.Serve(conn)
		} else {

			log.Fatalln("[MQTT ] Error:\n", err)
		}
	}
}

func ListenAndServeMQTTTLS(config *tls.Config) error {

	log.Println("[MQTTS] MQTT (with TLS) Server at \":8883\".")

	listener, err := tls.Listen("tcp", ":8883", config)
	if err != nil {
		log.Fatalln("[MQTTS] Error:\n", err)
	}

	for {

		conn, err := listener.Accept()
		if err == nil {

			go mqttServer.Serve(conn)
		} else {

			log.Fatalln("[MQTTS] Error:\n", err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

type MQTTResponse struct {
	status int
	header http.Header
}

func (resp *MQTTResponse) Header() http.Header {
	return resp.header
}

func (resp *MQTTResponse) Write(data []byte) (int, error) {
	return len(data), nil
}

func (resp *MQTTResponse) WriteHeader(statusCode int) {
	resp.status = statusCode
}

type MQTTHandler struct{}

func (h *MQTTHandler) Connect(conn *mqtt.Connection, username, password string) error {
	log.Printf("[MQTT ] (%s) Connect: %s, %s\n", conn.ClientID, username, password)
	return nil
}

func (h *MQTTHandler) Disconnect(conn *mqtt.Connection) {
	log.Printf("[MQTT ] (%s) Disconnect.\n", conn.ClientID)
}

func (h *MQTTHandler) Publish(conn *mqtt.Connection, msg *mqtt.Message) error {
	if conn != nil {
		log.Printf("[MQTT ] (%s) Published \"%s\" [%d].\n", conn.ClientID, msg.Topic, len(msg.Buf))

		body := tools.ClosingBuffer{bytes.NewBuffer(msg.Buf)}
		rurl, _ := url.Parse(msg.Topic)
		req := http.Request{
			Method: "PUBLISH",
			URL:    rurl,
			Header: http.Header{
				"X-Tag": []string{"MQTT "},
			},
			Body:          &body,
			ContentLength: int64(len(msg.Buf)),
			RemoteAddr:    conn.ClientID,
			RequestURI:    msg.Topic,
		}
		resp := MQTTResponse{
			status: 200,
			header: make(http.Header),
		}
		Serve(&resp, &req)
	}
	return nil
}

func (h *MQTTHandler) Subscribe(conn *mqtt.Connection, topic string, qos byte) error {
	log.Printf("[MQTT ] (%s) Subscribe \"%s\".\n", conn.ClientID, topic)
	return nil
}
