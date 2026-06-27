package main

import (
	"log"
	"net/http"

	"github.com/minaamim/go-chat-server/internal/chat"
	"github.com/minaamim/go-chat-server/internal/http/server"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	router := server.NewRouter(hub)
	log.Println("server started :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
