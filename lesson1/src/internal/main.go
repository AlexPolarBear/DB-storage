package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	port = flag.Int("port", 5000, "port to run the server on")
)

func main() {
	// listenAddr := os.Getenv("LISTEN_ADDR")
	// if len(listenAddr) == 0 {
	// 	listenAddr = "localhost:5000"
	// }
	flag.Parse()

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
// and not to output an error
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
// Creates a new file.bin and writes the receives data there
func replaceHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("file.bin")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("POST successfully")
	w.WriteHeader(200)
}

// GET "/get"
// Reads data from the file.bin and outputs to the screen
func getHandler(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("file.bin")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	data := make([]byte, 1024)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		w.Write(data[:n])
	}

	log.Println("GET successfully")
}

// GET /yummy
// You asked for it
func yummyHandler(w http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadFile("sid.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(buf)

	log.Println("HASH BROWNS")
}
