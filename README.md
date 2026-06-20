## Go Chat Server

A real-time chat server built with Go, Chi, and WebSocket to explore Go concurrency patterns such as goroutines and channels.

### Features

* WebSocket-based real-time communication
* Multi-client message broadcasting
* Concurrent connection handling with goroutines
* Channel-based message delivery
* HTTP routing with Chi

### Tech Stack

* Go
* Chi
* WebSocket
* Goroutines
* Channels

### Project Structure

```
go-chat-server/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── chat/
│   └── http/
│       ├── router.go
│       └── handler/
├── web/
├── go.mod
└── README.md
```

### Getting Started

Install dependencies

`go mod tidy`

Run the server

`go run ./cmd/server`

The server will start on:

`http://localhost:8080`

Health Check

`curl http://localhost:8080/health`

Expected response:

`ok`

### Learning Goals

This project was created to better understand:

* Goroutines
* Channels
* Interfaces
* WebSocket communication
* Concurrent programming in Go
* Dependency injection and application structure

### Roadmap

* Project setup
* Chi router
* Health check endpoint
* WebSocket connection
* Chat hub
* Client management
* Message broadcasting
* Chat rooms
* Redis Pub/Sub
* Graceful shutdown
* Authentication

License

MIT
