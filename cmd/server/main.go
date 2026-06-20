package main

import (
	"log"
	"net/http"

	"github.com/minaamim/go-chat-server/internal/http/server"
)

func main() {
	router := server.NewRouter()

	log.Println("server started :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
