package chat

import "log"

type Hub struct {
	Clients map[*Client]bool

	// 새로운 사용자가 접속했을 때 사용하는 채널
	Register chan *Client
	// 사용자가 연결을 끊었을 때 사용하는 채널
	Unregister chan *Client
	// 채팅 메세지를 모든 사용자에게 전달하기 위한 채널
	Broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("client registered, total=%d", len(h.Clients))

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("client unregistered, total=%d", len(h.Clients))
			}

		case message := <-h.Broadcast:
			log.Printf("broadcasting: %s to %d clients", message, len(h.Clients))

			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
