package handler

import (
	"net/http"
	"strings"

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

		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" {
			name = "noname"
		}
		client := chat.NewClient(hub, conn, name)

		if !hub.Register(client) {
			_ = conn.Close()
			return
		}

		go client.WritePump()
		client.ReadPump()
	}
}
