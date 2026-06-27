package chat

import "github.com/gorilla/websocket"

const sendBufferSize = 256

type Client struct {
	hub *Hub
	// 브라우저와 연결된 실제 Websocket 연결
	conn *websocket.Conn
	// 사용자에게 보낼 메세지 리스트
	send chan []byte
	name string
}

func NewClient(hub *Hub, conn *websocket.Conn, name string) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, sendBufferSize),
		name: name,
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	for message := range c.send {
		err := c.conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {
			return
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.Broadcast(c, msg)
	}
}
