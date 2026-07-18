package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minaamim/go-chat-server/internal/chat"
	"github.com/minaamim/go-chat-server/internal/http/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	hub := chat.NewHub()
	go hub.Run(ctx)

	router := server.NewRouter(hub)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("server started :8080")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}

	select {
	case <-hub.Done():
	case <-shutdownCtx.Done():
		log.Printf("hub shutdown timed out: %v", shutdownCtx.Err())
	}
}
