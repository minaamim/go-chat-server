package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/minaamim/go-chat-server/internal/chat"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocket(hub *chat.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		client := chat.NewClient(hub, conn)

		hub.Register(client)

		go client.WritePump()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				hub.Unregister(client)
				break
			}
			hub.Broadcast(msg)
		}
	}
}
