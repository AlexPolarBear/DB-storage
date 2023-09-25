package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	transaction "../lesson2/transaction/manager.go"
)

var (
	port  = flag.Int("port", 5000, "port to run the server on")
	queue = make(chan []byte)
	data  []byte
)

// Creates connection to http server.
func main() {
	flag.Parse()

	transaction.NewTransactionManager(queue)

	mux := http.NewServeMux()
	mux.HandleFunc("/", docs)
	mux.HandleFunc("/replace", replaceHandler)
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/yummy", yummyHandler)

	log.Printf("Connection on: http://localhost:%d \n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

// GET "/"
// Easy api to show what methods are available,
// and not to output an error.
func docs(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	w.Write([]byte("Hello! This is main page.\nThere are several methods here:\n"))
	w.Write([]byte("\t\xF0\x9F\x8D\x95 /get -- GET body from a file\n"))
	w.Write([]byte("\t\xF0\x9F\x8D\x94 /replace -- POST body to a file\n"))
	w.Write([]byte("\t\xF0\x9F\x8D\xA0 /yummy -- GET hash browns\n"))

	log.Println("HOME page")
}

// POST "/replace"
// Creates a new file.bin and writes the receives data there.
func replaceHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	queue <- body
	data = body

	log.Println("POST successfully")
	w.WriteHeader(200)
}

// GET "/get"
// Reads data from the file.bin and outputs to the screen.
func getHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GET successfully")
}

// GET /yummy
// You asked for it.
func yummyHandler(w http.ResponseWriter, req *http.Request) {
	buf, err := os.ReadFile("sid.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(buf)

	log.Println("HASH BROWNS")
}
