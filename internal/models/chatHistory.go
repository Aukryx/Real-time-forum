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
	Type     string               `json:"type"`
	User1ID  int                  `json:"user1_id"`
	User2ID  int                  `json:"user2_id"`
	Messages []ChatHistoryMessage `json:"messages"`
}
