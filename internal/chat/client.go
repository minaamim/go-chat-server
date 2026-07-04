package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	sendBufferSize = 256
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
)

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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(
				websocket.TextMessage,
				message,
			); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.Broadcast(c, msg)
	}
}
