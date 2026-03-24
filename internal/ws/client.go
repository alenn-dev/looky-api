package ws

import (
	"log"

	"github.com/fasthttp/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(1, msg); err != nil {
			log.Println("websocket write error:", err)
			break
		}
	}
}

func (c *Client) ReadPump(customerID string) {
	defer func() {
		H.UnRegister(customerID)
		close(c.send)
		c.conn.Close()
	}()

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}
