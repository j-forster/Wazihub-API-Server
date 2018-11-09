package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/j-forster/Wazihub-API/mqtt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	req.Header.Set("X-Secure", "false")
	req.Header.Set("X-Tag", "HTTP ")
	serveHTTP(resp, req)
}

func ServeHTTPS(resp http.ResponseWriter, req *http.Request) {

	req.Header.Set("X-Secure", "true")
	req.Header.Set("X-Tag", "HTTPS")
	serveHTTP(resp, req)
}

////////////////////

type wsWrapper struct {
	conn   *websocket.Conn
	wc     io.WriteCloser
	head   mqtt.FixedHeader
	remain int
}

func (w *wsWrapper) Close() error {
	return w.conn.Close()
}

func (w *wsWrapper) Write(data []byte) (int, error) {

	if w.remain == 0 {
		w.wc, _ = w.conn.NextWriter(websocket.BinaryMessage)
		l, _ := w.head.ReadMessage(data)
		w.remain = l + w.head.Length
	}

	w.remain -= len(data)
	_, err := w.wc.Write(data)
	if w.remain == 0 {
		w.wc.Close()
	}
	return len(data), err
}

////////////////////

func serveHTTP(resp http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Upgrade") != "websocket" {

		Serve(resp, req) // see main.go
	} else {

		proto := req.Header.Get("Sec-WebSocket-Protocol")
		if proto != "mqttv3.1" {
			http.Error(resp, "Requires WebSocket Protocol Header 'mqttv3.1'.", http.StatusBadRequest)
			return
		}

		responseHeader := make(http.Header)
		responseHeader.Set("Sec-WebSocket-Protocol", "mqttv3.1")

		conn, err := upgrader.Upgrade(resp, req, responseHeader)
		if err != nil {
			log.Printf("[%s] (%s) WebSocket Upgrade Failed\n %v", req.Header.Get("X-Tag"), req.RemoteAddr, err)
			return
		}

		var tag string
		if req.Header.Get("X-Secure") == "true" {
			tag = "WSS  "
		} else {
			tag = "WS   "
		}

		wrapper := wsWrapper{conn: conn}
		mqttConn := mqtt.NewConnection(&wrapper, &wrapper, mqttServer)

		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[%s] (%s) WebSocket Read Error\n %v", tag, conn.RemoteAddr().String(), err)
				//mqttConn.Close()
				conn.Close() // obsolete
				return
			}

			if messageType != websocket.BinaryMessage {
				log.Printf("[%s] (%s) WebSocket Error:\n Unexpected TEXT message.", tag, conn.RemoteAddr().String())
				mqttConn.Close()
				conn.Close() // obsolete
				return
			}

			mqttConn.ReadMessage(msg)
			/*
				if err := conn.WriteMessage(messageType, msg); err != nil {
					log.Printf("[%s] (%s) WebSocket Write Error\n %v", tag, conn.RemoteAddr().String(), err)
					return
				}
			*/

		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func ListenAndServeHTTP() {

	log.Println("[HTTP ] HTTP Server at \":80\". Use \"http://\".")
	log.Println("[WS   ] MQTT via WebSocket Server at \":80\". Use \"ws://\".")
	err := http.ListenAndServe(":80", http.HandlerFunc(ServeHTTP))
	if err != nil {
		log.Println("[HTTP ] Error:")
		log.Fatalln(err)
	}
}

func ListenAndServeHTTPS(cfg *tls.Config) {

	srv := &http.Server{
		Addr:         ":443",
		Handler:      http.HandlerFunc(ServeHTTPS),
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	log.Println("[HTTPS] HTTPS Server at \":443\". Use \"https://\".")
	log.Println("[WSS  ] MQTT via WebSocket Server at \":443\".  Use \"wss://\".")
	err := srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalln("[HTTPS] Error:\n", err)
	}
}
