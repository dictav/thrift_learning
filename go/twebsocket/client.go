package twebsocket

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxID int

// Client struct
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *string
	doneCh chan bool
}

// NewClient creates new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxID++
	ch := make(chan *string, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxID, ws, server, ch, doneCh}
}

// Conn function
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

// Write function
func (c *Client) Write(msg *string) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected", c.id)
		c.server.Err(err)
	}
}

// Done function
func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", *msg)
			err := websocket.Message.Send(c.ws, *msg)
			if err != nil {
				c.server.Err(err)
			}

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var str string
			err := websocket.Message.Receive(c.ws, &str)
			log.Println("Receive", str)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				ret, err := c.server.Process(&str)
				if err == nil {
					c.Write(&ret)
				}
			}
		}
	}
}
