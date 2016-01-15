package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/dictav/thrift_learning/go/gen-go/tutorial"
	"github.com/dictav/thrift_learning/go/thttp"
	"github.com/dictav/thrift_learning/go/twebsocket"
)

func main() {
	h := NewCalculatorHandler()
	processor := tutorial.NewCalculatorProcessor(h)

	mux := http.NewServeMux()
	// http server
	mux.HandleFunc("/calc", thttp.NewThriftHandler(processor))

	// websocket server
	server := twebsocket.NewServer("/ws_calc", processor)
	go server.Listen()

	// serve static files
	mux.Handle("/", http.FileServer(http.Dir("../browser")))

	// start negroni
	n := negroni.New(NewAuth())
	n.UseHandler(mux)
	n.Run(":9090")
}
