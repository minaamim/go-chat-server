package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Println("client connected")
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("client disconnected")
			break
		}
		log.Printf("received: %s\n", msg)
	}
}
