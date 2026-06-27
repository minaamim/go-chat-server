package chat

import "github.com/gorilla/websocket"

type Client struct {
	Hub *Hub
	// 브라우저와 연결된 실제 Websocket 연결
	Conn *websocket.Conn
	// 사용자에게 보낼 메세지 리스트
	Send chan []byte
}

func (c *Client) WritePump() {
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	for message := range c.Send {
		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {
			return
		}
	}
}
