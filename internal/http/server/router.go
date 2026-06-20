package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/minaamim/go-chat-server/internal/http/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)
	r.Get("/ws", handler.WebSocket)
	return r
}
