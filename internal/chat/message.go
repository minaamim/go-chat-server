package chat

type Message struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
