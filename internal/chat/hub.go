package chat

import "log"

type Hub struct {
	clients map[*Client]bool

	// 새로운 사용자가 접속했을 때 사용하는 채널
	register chan *Client
	// 사용자가 연결을 끊었을 때 사용하는 채널
	unregister chan *Client
	// 채팅 메세지를 모든 사용자에게 전달하기 위한 채널
	broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("client registered, total=%d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("client unregistered, total=%d", len(h.clients))
			}

		case message := <-h.broadcast:
			log.Printf("broadcasting: %s to %d clients", message, len(h.clients))

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
