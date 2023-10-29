package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type TransactionStruct struct {
	Source  string // ваша фамилия
	Id      uint64 // возрастающий счетчик
	Payload string // транзакция (json patch)
}

var (
	// host      = flag.String("host", "127.0.0.1", "host to run server on")
	port      = flag.Int("port", 5000, "port to run the server on")
	source    = "Alex-Polar-Bear"
	counterID = uint64(0)
	vclock    = make(map[string]uint64)
	queue     = make(chan *TransactionStruct)
	snapshot  = "{ \"baz\": \"qux\", \"foo\": \"bar\" }"
	peers     = []string{"127.0.0.1:5001"}
	wsConn    = make([]*websocket.Conn, 0)
	wal       = make([]*TransactionStruct, 0) // журнал транзакций

	//go:embed config/index.html
	test string
)

func main() {
	flag.Parse()

	go ManagerTransaction()

	for _, peer := range peers {
		go webSocket(peer)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/vclock", vclockHandler)
	mux.HandleFunc("/replace", replaceHandler)
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/ws", wsHandler)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)

	//log.Printf("Connection on: http://%s:%d \n", *host, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(test))
	if err != nil {
		log.Fatal(err)
	}
}

func vclockHandler(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(vclock)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func replaceHandler(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	counterID += 1

	tm := &TransactionStruct{
		Source:  source,
		Id:      counterID,
		Payload: string(data),
	}

	queue <- tm
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(snapshot))
	if err != nil {
		log.Fatal(err)
	}
}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{" (J o_o)J "},
	})

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Reecieved new connection ", "from", req.RemoteAddr)
	wsConn = append(wsConn, conn)
}

func webSocket(peer string) {
	var conn *websocket.Conn
	var err error
	ctx := context.Background()
	for {
		conn, _, err = websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", peer), nil)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(1 * time.Minute)
			continue
		}
		break
	}

	for {
		var tr TransactionStruct
		err = wsjson.Read(ctx, conn, &tr)
		if errors.Is(err, io.EOF) {
			slog.Info("Peer disconnected", "peer", peer)
			break
		}

		slog.Info("Received", "transaction", tr)
		if err != nil {
			log.Fatal(err)
		}

		queue <- &tr
	}
}

func ManagerTransaction() {
	for {
		tr := <-queue

		if vclock[tr.Source] > tr.Id {
			continue
		}
		vclock[tr.Source] = tr.Id + 1
		wal = append(wal, tr)
		patch, err := jsonpatch.DecodePatch([]byte(tr.Payload))
		if err != nil {
			panic(err)
		}

		snapBytes, err := patch.Apply([]byte(snapshot))
		if err != nil {
			panic(err)
		}

		snapshot = string(snapBytes)

		slog.Info("Sending transaction to peers", "transaction", *tr)
		for _, conn := range wsConn {
			err = wsjson.Write(context.Background(), conn, tr)
			if err != nil {
				panic(err)
			}
		}
	}
}
