package internal

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	host = flag.String("host", "127.0.0.1", "host to run server on")
	port = flag.Int("port", 5000, "port to run the server on")
	// queue = make(chan []byte)
	// data  []byte
)

func main() {
	flag.Parse()

	// transaction.NewTransactionManager(queue)

	mux := http.NewServeMux()
	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/vclock", vclockHandler)
	mux.HandleFunc("/replace", replaceHandler)
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/ws", wsHandler)

	log.Printf("Connection on: http://%s:%d \n", *host, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	panic("implement me")
}

func vclockHandler(w http.ResponseWriter, req *http.Request) {
	panic("implement me")
}

func replaceHandler(w http.ResponseWriter, req *http.Request) {
	panic("implement me")
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	panic("implement me")
}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	panic("implement me")
}
