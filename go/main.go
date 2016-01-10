package main

import (
	"log"
	"net/http"

	"github.com/dictav/thrift_learning/go/gen-go/tutorial"
	"github.com/dictav/thrift_learning/go/thttp"
)

func main() {
	h := NewCalculatorHandler()
	processor := tutorial.NewCalculatorProcessor(h)

	http.HandleFunc("/calc", thttp.NewThriftHandler(processor))
	http.Handle("/", http.FileServer(http.Dir("../browser")))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
