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

### Architecture

```

                         +----------------------+
                         |         Hub          |
                         |----------------------|
                         | Clients              |
                         | Register             |
                         | Unregister           |
                         | Broadcast            |
                         +----------+-----------+
                                    ^
                                    |
             Register / Broadcast / Unregister
                                    |
        +---------------------------+---------------------------+
        |                           |                           |
        |                           |                           |
+-------------------+     +-------------------+     +-------------------+
|      Client       |     |      Client       |     |      Client       |
|-------------------|     |-------------------|     |-------------------|
| Conn              |     | Conn              |     | Conn              |
| Send (buffered)   |     | Send (buffered)   |     | Send (buffered)   |
| Name              |     | Name              |     | Name              |
+---------+---------+     +---------+---------+     +---------+---------+
|                           |                           |
+------+-------+            +------+-------+            +------+-------+
| ReadPump()   |            | ReadPump()   |            | ReadPump()   |
| WritePump()  |            | WritePump()  |            | WritePump()  |
+------+-------+            +------+-------+            +------+-------+
|                           |                           |
+------------ WebSocket Connections --------------------+

```

### Message Flow

```
        Browser A                 Browser B
             │                        │
        WebSocket                WebSocket
             │                        │
     +-------▼--------+      +--------▼-------+
     | Read / Write   |      | Read / Write   |
     |     Pump       |      |     Pump       |
     +-------+--------+      +--------+-------+
             │                        │
             └──────────┬─────────────┘
                        │
                Register / Broadcast
                        │
                 +------+------+
                 |     Hub     |
                 +------+------+
                        │
                Broadcast Message
                        │
             ┌──────────┴──────────┐
             ▼                     ▼
        Client.Send          Client.Send
```

### Project Structure

```
go-chat-server/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── chat/
│   │   ├── client.go            # WebSocket client read/write pumps
│   │   ├── hub.go               # Client registry and message broadcasting
│   │   ├── hub_test.go          # Hub and message behavior tests
│   │   └── message.go           # Chat and system message payloads
│   └── http/
│       ├── handler/
│       │   ├── health.go        # Health check handler
│       │   └── websocket.go     # WebSocket upgrade handler
│       └── server/
│           └── router.go        # Chi router setup
├── web/
│   └── index.html               # Browser test client
├── go.mod
├── go.sum
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

### Key Concepts

This project was created to better understand:

* Goroutines
* Channels
* Interfaces
* WebSocket communication
* Concurrent programming in Go
* Dependency injection and application structure

### Hub Concurrency Model

The chat hub uses a single owner goroutine to manage connected clients.

`Hub.Run()` is the only place that reads or writes the `clients` map. Other goroutines do not touch the map directly. They communicate with the hub through channels:

* `Register` sends a client to the hub when a WebSocket connection is created.
* `Unregister` sends a client to the hub when a read or write pump exits.
* `Broadcast` sends a chat event to the hub so it can fan out one payload to every connected client.

Because one goroutine owns the map, the current design does not need `sync.RWMutex`. A mutex would become useful if multiple goroutines started reading or writing `clients` directly, for example through a separate user-count API that accessed the map outside `Hub.Run()`.

Each client has a buffered `send` channel. The hub writes outbound payloads into that channel, and the client's write pump is responsible for writing them to the WebSocket connection. If a client's send buffer is full, the hub treats that client as slow or disconnected, closes its send channel, and removes it from the map. This prevents one slow client from blocking broadcasts for everyone else.

WebSocket connection lifetime is split between two pumps:

* `ReadPump` reads incoming messages, broadcasts them through the hub, and unregisters the client when reading fails.
* `WritePump` writes outbound messages and periodic ping frames, and unregisters the client when writing fails.

The server uses ping/pong deadlines to detect stale WebSocket connections. The write pump sends ping frames on an interval, and the read pump extends the read deadline whenever a pong frame is received.

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

### License

MIT
