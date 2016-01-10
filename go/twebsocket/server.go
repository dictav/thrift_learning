package twebsocket

import (
	"bytes"
	"log"
	"net/http"

	"git.apache.org/thrift.git/lib/go/thrift"
	"golang.org/x/net/websocket"
)

// Server struct
type Server struct {
	pattern   string
	processor thrift.TProcessor
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	doneCh    chan bool
	errCh     chan error
}

// NewServer creates new thrift websocket server.
func NewServer(pattern string, processor thrift.TProcessor) *Server {
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		processor,
		clients,
		addCh,
		delCh,
		doneCh,
		errCh,
	}
}

// Add function
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

// Del function
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// Process function
func (s *Server) Process(str *string) (string, error) {
	buffer := new(bytes.Buffer)
	_, err := buffer.WriteString(*str)
	if err != nil {
		s.Err(err)
		return "", err
	}

	t1 := &thrift.TMemoryBuffer{Buffer: buffer}
	in := thrift.NewTJSONProtocol(t1)
	t2 := thrift.NewTMemoryBufferLen(1024)
	//defer t2.Close()
	out := thrift.NewTJSONProtocol(t2)

	s.processor.Process(in, out)
	return t2.String(), nil
}

// Done function
func (s *Server) Done() {
	s.doneCh <- true
}

// Err function
func (s *Server) Err(err error) {
	s.errCh <- err
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")

			// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
