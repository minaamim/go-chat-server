package chat

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHubBroadcastsChatMessagesWithSenderName(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	alice := NewClient(hub, nil, "alice")
	bob := NewClient(hub, nil, "bob")

	hub.Register(alice)
	assertMessage(t, receiveMessage(t, alice), Message{
		Type:    "system",
		Content: "alice joined",
	})

	hub.Register(bob)
	assertMessage(t, receiveMessage(t, alice), Message{
		Type:    "system",
		Content: "bob joined",
	})
	assertMessage(t, receiveMessage(t, bob), Message{
		Type:    "system",
		Content: "bob joined",
	})

	hub.Broadcast(alice, []byte("hello"))

	want := Message{
		Type:    "chat",
		Name:    "alice",
		Content: "hello",
	}
	assertMessage(t, receiveMessage(t, alice), want)
	assertMessage(t, receiveMessage(t, bob), want)
}

func TestHubBroadcastsLeaveMessagesToRemainingClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	alice := NewClient(hub, nil, "alice")
	bob := NewClient(hub, nil, "bob")

	hub.Register(alice)
	receiveMessage(t, alice)

	hub.Register(bob)
	receiveMessage(t, alice)
	receiveMessage(t, bob)

	hub.Unregister(bob)

	assertMessage(t, receiveMessage(t, alice), Message{
		Type:    "system",
		Content: "bob left",
	})
}

func TestHubRemovesClientsWithFullSendBuffer(t *testing.T) {
	hub := NewHub()
	client := NewClient(hub, nil, "slow")
	hub.clients[client] = true

	for i := 0; i < sendBufferSize; i++ {
		client.send <- []byte("queued")
	}

	hub.broadcastPayload([]byte(`{"type":"system","content":"overflow"}`))

	if _, ok := hub.clients[client]; ok {
		t.Fatal("expected slow client to be removed")
	}
}

func TestMessageJSONFields(t *testing.T) {
	payload, err := json.Marshal(Message{
		Type:    "chat",
		Name:    "alice",
		Content: "hello",
	})
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	const want = `{"type":"chat","name":"alice","content":"hello"}`
	if string(payload) != want {
		t.Fatalf("message json = %s, want %s", payload, want)
	}
}

func receiveMessage(t *testing.T, client *Client) Message {
	t.Helper()

	select {
	case payload := <-client.send:
		var message Message
		if err := json.Unmarshal(payload, &message); err != nil {
			t.Fatalf("unmarshal message: %v", err)
		}
		return message
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for message")
		return Message{}
	}
}

func assertMessage(t *testing.T, got, want Message) {
	t.Helper()

	if got != want {
		t.Fatalf("message = %+v, want %+v", got, want)
	}
}
