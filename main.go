package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/j-forster/Wazihub-API/mqtt"
	"github.com/j-forster/Wazihub-API/tools"
)

/*
var db *mongo.Client
var collection *mongo.Collection
*/

var db *mgo.Session
var collection *mgo.Collection

var static http.Handler

func main() {
	// Remove date and time from logs
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		interactive()
		return
	}

	tlsCert := flag.String("crt", "", "TLS Cert File (.crt)")
	tlsKey := flag.String("key", "", "TLS Key File (.key)")

	www := flag.String("www", "/var/www", "HTTP files root")

	nodb := flag.Bool("no-db", false, "Disable MongoDB")

	dbAddr := flag.String("db", "localhost:27017", "MongoDB address.")

	upstream, _ := getConfig("upstream")
	upstreamAddr := flag.String("upstream", upstream, "Upstream server address.")

	flag.Parse()

	////////////////////

	log.Println("WaziHub API Server")
	log.Println("--------------------")

	////////////////////

	if *www != "" {
		static = http.FileServer(http.Dir(*www))
	}

	////////////////////

	if !*nodb {

		log.Printf("[DB   ] Dialing MongoDB at %q...\n", *dbAddr)

		var err error
		db, err = mgo.Dial("mongodb://" + *dbAddr + "/?connect=direct")
		if err != nil {
			db = nil
			log.Println("[DB   ] MongoDB client error:\n", err)
		} else {

			collection = db.DB("Wazihub").C("values")
		}
	}

	////////////////////

	if *upstreamAddr != "" {
		go Upstream(*upstreamAddr)
	}

	////////////////////

	if *tlsCert != "" && *tlsKey != "" {

		cert, err := ioutil.ReadFile(*tlsCert)
		if err != nil {
			log.Println("Error reading", *tlsCert)
			log.Fatalln(err)
		}

		key, err := ioutil.ReadFile(*tlsKey)
		if err != nil {
			log.Println("Error reading", *tlsKey)
			log.Fatalln(err)
		}

		pair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			log.Println("TLS/SSL 'X509KeyPair' Error")
			log.Fatalln(err)
		}

		cfg := &tls.Config{Certificates: []tls.Certificate{pair}}

		ListenAndServeHTTPS(cfg)
		ListenAndServeMQTTTLS(cfg)
	}

	ListenAndServerMQTT()
	ListenAndServeHTTP() // will block
}

///////////////////////////////////////////////////////////////////////////////

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (resp *ResponseWriter) WriteHeader(statusCode int) {
	resp.status = statusCode
	resp.ResponseWriter.WriteHeader(statusCode)
}

////////////////////

func Serve(resp http.ResponseWriter, req *http.Request) {
	wrapper := ResponseWriter{resp, 200}

	if static != nil {
		if strings.HasPrefix(req.RequestURI, "/www/") {
			req.RequestURI = req.RequestURI[4:]
			req.URL.Path = req.URL.Path[4:]
			static.ServeHTTP(&wrapper, req)

			log.Printf("[WWW  ] (%s) %d %s \"/www%s\"\n",
				req.RemoteAddr,
				wrapper.status,
				req.Method,
				req.RequestURI)
			return
		}
	}

	if req.Method == http.MethodPut || req.Method == http.MethodPost {

		body, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			http.Error(resp, "400 Bad Request", http.StatusBadRequest)
			return
		}
		req.Body = &tools.ClosingBuffer{bytes.NewBuffer(body)}
	}

	if req.Method == "PUBLISH" {
		req.Method = http.MethodPost
		router.ServeHTTP(&wrapper, req)
		req.Method = "PUBLISH"
	} else {
		router.ServeHTTP(&wrapper, req)
	}

	/*
		log.Printf("[%s] (%s) %d %s \"%s\"\n",
			req.Header.Get("X-Tag"),
			req.RemoteAddr,
			wrapper.status,
			req.Method,
			req.RequestURI)
	*/

	if cbuf, ok := req.Body.(*tools.ClosingBuffer); ok {
		// log.Printf("[DEBUG] Body: %s\n", cbuf.Bytes())
		msg := &mqtt.Message{
			QoS:   0,
			Topic: req.RequestURI[1:],
			Data:  cbuf.Bytes(),
		}

		// if wrapper.status >= 200 && wrapper.status < 300 {
		if req.Method == http.MethodPut || req.Method == http.MethodPost {
			mqttServer.Publish(nil, msg)
		}
		// }
	}

}
