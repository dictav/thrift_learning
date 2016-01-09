package main

import (
	"bytes"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"net/http"

	"github.com/dictav/thrift_learning/go/gen-go/tutorial"
)

func barHandler(w http.ResponseWriter, r *http.Request) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r.Body)
	inTransport := &thrift.TMemoryBuffer{Buffer: buffer}
	inProtocol := thrift.NewTJSONProtocol(inTransport)
	outTransport := thrift.NewTMemoryBufferLen(1024)
	outProtocol := thrift.NewTJSONProtocol(outTransport)

	handler := NewCalculatorHandler()
	processor := tutorial.NewCalculatorProcessor(handler)
	processor.Process(inProtocol, outProtocol)

	fmt.Printf("%q: out=%q\n", r.URL.Path, outTransport)
	fmt.Fprint(w, outTransport)
}

func main() {
	http.HandleFunc("/calc", barHandler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
