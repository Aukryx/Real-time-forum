package models

// ChatHistoryMessage represents a single message in the chat history
type ChatHistoryMessage struct {
	Type      string `json:"type"`
	SenderID  int    `json:"sender_id"`
	Sender    string `json:"sender"` // Username of the sender
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"` // Creation time as string
	Read      bool   `json:"read"`      // Whether the message has been read
}

// ChatHistory represents the full history of messages between two users
type ChatHistory struct {
	Type      string               `json:"type"`
	User1Name string               `json:"user1name"`
	User2Name string               `json:"user2name"`
	Messages  []ChatHistoryMessage `json:"messages"`
}
