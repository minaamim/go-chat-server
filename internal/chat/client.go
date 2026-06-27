package chat

import "github.com/gorilla/websocket"

const sendBufferSize = 256

type Client struct {
	hub *Hub
	// 브라우저와 연결된 실제 Websocket 연결
	conn *websocket.Conn
	// 사용자에게 보낼 메세지 리스트
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, sendBufferSize),
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
