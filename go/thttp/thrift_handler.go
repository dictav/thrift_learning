package thttp

import (
	"bytes"
	"fmt"
	"net/http"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// NewThriftHandler retuns http.Handler
func NewThriftHandler(processor thrift.TProcessor) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != `POST` {
			http.Error(w, "Thrift Handler requires POST access", http.StatusInternalServerError)
			return
		}

		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}

		inTransport := &thrift.TMemoryBuffer{Buffer: buffer}
		inProtocol := thrift.NewTJSONProtocol(inTransport)
		outTransport := thrift.NewTMemoryBufferLen(1024)
		outProtocol := thrift.NewTJSONProtocol(outTransport)

		processor.Process(inProtocol, outProtocol)

		fmt.Fprint(w, outTransport)
	}
}
