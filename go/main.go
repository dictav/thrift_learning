package main

import (
	"log"
	"net/http"

	"github.com/dictav/thrift_learning/go/gen-go/tutorial"
	"github.com/dictav/thrift_learning/go/thttp"
	"github.com/dictav/thrift_learning/go/twebsocket"
)

func main() {
	h := NewCalculatorHandler()
	processor := tutorial.NewCalculatorProcessor(h)

	// http server
	http.HandleFunc("/calc", thttp.NewThriftHandler(processor))

	// websocket server
	server := twebsocket.NewServer("/ws_calc", processor)
	go server.Listen()

	// serve static files
	http.Handle("/", http.FileServer(http.Dir("../browser")))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
