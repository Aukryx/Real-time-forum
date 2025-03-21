package models

// Struct that will store the content of the json sent by the javascript
type PrivateMessage struct {
	Type     string   `json:"type"`
	Sender   string   `json:"sender"`
	Receiver string   `json:"receiver"`
	Message  string   `json:"message"`
	UserList []string `json:"user_list,omitempty"`
}
