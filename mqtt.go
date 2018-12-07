package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/j-forster/Wazihub-API/mqtt"
	"github.com/j-forster/Wazihub-API/tools"
)

type MQTTServer struct {
	mqtt.Server
}

var connectCounter, disconnectCounter int

func (server *MQTTServer) Publish(client *mqtt.Client, msg *mqtt.Message) error {

	if client != nil {

		log.Printf("[MQTT ] Publish %q %q QoS:%d len:%d\n", client.Id, msg.Topic, msg.QoS, len(msg.Data))

		body := tools.ClosingBuffer{bytes.NewBuffer(msg.Data)}
		rurl, _ := url.Parse(msg.Topic)
		req := http.Request{
			Method: "PUBLISH",
			URL:    rurl,
			Header: http.Header{
				"X-Tag": []string{"MQTT "},
			},
			Body:          &body,
			ContentLength: int64(len(msg.Data)),
			RemoteAddr:    client.Id,
			RequestURI:    msg.Topic,
		}
		resp := MQTTResponse{
			status: 200,
			header: make(http.Header),
		}

		Serve(&resp, &req)
		// if resp.status >= 200 && resp.status < 300 {
		server.Server.Publish(client, msg)
		// }
	} else {

		server.Server.Publish(client, msg)
	}

	if upstream != nil {

		upstream <- msg
	}

	return nil
}

func (server *MQTTServer) Connect(client *mqtt.Client, auth *mqtt.ConnectAuth) byte {
	log.Printf("[MQTT ] Connect %q %+v\n", client.Id, auth)
	return mqtt.CodeAccepted
}

func (server *MQTTServer) Disconnect(client *mqtt.Client, err error) {
	log.Printf("[MQTT ] Disonnect %q %v\n", client.Id, err)
}

func (server *MQTTServer) Subscribe(recv mqtt.Reciever, topic string, qos byte) *mqtt.Subscription {

	if client, ok := recv.(*mqtt.Client); ok {
		log.Printf("[MQTT ] Subscribe %q %q QoS:%d\n", client.Id, topic, qos)
	}
	return server.Server.Subscribe(recv, topic, qos)
}

func (server *MQTTServer) Unsubscribe(subs *mqtt.Subscription) {

	if client, ok := subs.Recv.(*mqtt.Client); ok {
		log.Printf("[MQTT ] Unsubscribe %q %q QoS:%d\n", client.Id, subs.Topic.FullName(), subs.QoS)
	}
	server.Server.Unsubscribe(subs)
}

////////////////////////////////////////////////////////////////////////////////

var mqttServer = &MQTTServer{mqtt.NewServer()}

func ListenAndServerMQTT() {

	log.Println("[MQTT ] MQTT Server at \":1883\".")
	mqtt.ListenAndServe(":1883", mqttServer)
}

func ListenAndServeMQTTTLS(config *tls.Config) {

	log.Println("[MQTTS] MQTT (with TLS) Server at \":8883\".")
	mqtt.ListenAndServeTLS(":8883", config, mqttServer)
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
