package chat

import (
	"context"
	"encoding/json"
	"log"
)

type (
	Hub struct {
		clients map[*Client]bool

		// 새로운 사용자가 접속했을 때 사용하는 채널
		register chan *Client
		// 사용자가 연결을 끊었을 때 사용하는 채널
		unregister chan *Client
		// 채팅 메세지를 모든 사용자에게 전달하기 위한 채널
		broadcast chan broadcastMessage
		// Hub goroutine이 완전히 종료됐음을 알리는 채널
		done chan struct{}
	}

	broadcastMessage struct {
		sender  *Client
		content []byte
	}
)

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan broadcastMessage),
		done:       make(chan struct{}),
	}
}

func (h *Hub) Register(client *Client) bool {
	select {
	case h.register <- client:
		return true
	case <-h.done:
		return false
	}
}

func (h *Hub) Unregister(client *Client) {
	select {
	case h.unregister <- client:
	case <-h.done:
	}
}

func (h *Hub) Broadcast(sender *Client, content []byte) {
	select {
	case h.broadcast <- broadcastMessage{
		sender:  sender,
		content: content,
	}:
	case <-h.done:
	}
}

func (h *Hub) Done() <-chan struct{} {
	return h.done
}

func (h *Hub) broadcastPayload(payload []byte) {
	for client := range h.clients {
		select {
		case client.send <- payload:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) broadcastSystem(content string) {
	payload, err := json.Marshal(Message{
		Type:    "system",
		Content: content,
	})
	if err != nil {
		return
	}

	h.broadcastPayload(payload)
}

func (h *Hub) Run(ctx context.Context) {
	defer close(h.done)

	for {
		select {
		case <-ctx.Done():
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
			return

		case client := <-h.register:
			h.clients[client] = true
			log.Printf("client %s registered, total=%d", client.name, len(h.clients))
			h.broadcastSystem(client.name + " joined")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("client %s unregistered, total=%d", client.name, len(h.clients))
				h.broadcastSystem(client.name + " left")
			}

		case message := <-h.broadcast:
			payload, err := json.Marshal(Message{
				Type:    "chat",
				Name:    message.sender.name,
				Content: string(message.content),
			})
			if err != nil {
				continue
			}

			log.Printf("broadcasting: %s to %d clients", message.sender.name, len(h.clients))

			h.broadcastPayload(payload)
		}
	}
}
