package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/schigh/cancelme/downstream"
)

var (
	port int
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&port, "p", 8080, "port")
	flag.Parse()

	if port == 0 {
		log.Panic("port required")
	}

	mux := chi.NewRouter()
	mux.Get("/", handleFoo)

	log.Printf("listening on port %d\n", port)
	_ = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
}

func handleError(rw http.ResponseWriter, err error) {
	log.Printf("%v\n", err)

	rw.WriteHeader(500)
	_,_ = rw.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
}

func handleData(rw http.ResponseWriter, data map[string]interface{}) {
	log.Printf("%#v\n", data)
	b, _ := json.Marshal(&data)
	rw.WriteHeader(200)
	_, _ = rw.Write(b)
}

func handleFoo(rw http.ResponseWriter, r *http.Request) {
	foo := &downstream.Foo{}

	async, ok := r.URL.Query()["async"]
	if ok {
		foo.Async = async[0] == "true"
		log.Printf("async: %t\n", foo.Async)
	}

	asyncdies, ok := r.URL.Query()["asyncdies"]
	if ok {
		d, _ := strconv.Atoi(asyncdies[0])
		foo.AsyncDies = d
		log.Printf("asyncdies: %d\n", foo.AsyncDies)
	}

	timeout, ok := r.URL.Query()["timeout"]
	if ok {
		t, _ := strconv.ParseInt(timeout[0], 10, 64)
		foo.Timeout = time.Duration(t)
		log.Printf("timeout: %v\n", foo.Timeout)
	}

	depth, ok := r.URL.Query()["depth"]
	if ok {
		d, _ := strconv.Atoi(depth[0])
		foo.Depth = d
		log.Printf("depth: %d\n", foo.Depth)
	}

	cancels, ok := r.URL.Query()["cancels"]
	if ok {
		foo.Cancels = cancels[0] == "true"
		log.Printf("cancels: %t\n", foo.Cancels)
	}

	canceldepth, ok := r.URL.Query()["canceldepth"]
	if ok {
		d, _ := strconv.Atoi(canceldepth[0])
		foo.CancelDepth = d
		log.Printf("canceldepth: %d\n", foo.CancelDepth)
	}

	stagepause, ok := r.URL.Query()["stagepause"]
	if ok {
		t, _ := strconv.ParseInt(stagepause[0], 10, 64)
		foo.StagePause = time.Duration(t)
		log.Printf("stage pause: %v\n", foo.StagePause)
	}

	data, err := foo.Do(r.Context())
	if err != nil {
		handleError(rw, err)
		return
	}

	handleData(rw, data)
}
