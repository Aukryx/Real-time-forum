package models

import "time"

// Struct that will store the content of the json sent by the javascript
type PrivateMessage struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	Read      bool      `json:"read"`
}
