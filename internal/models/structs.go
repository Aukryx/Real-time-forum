package models

import (
	"time"
)

type User struct {
	ID        int
	UUID      string
	NickName  string
	Gender    string
	FirstName string
	LastName  string
	Password  string
	Email     string
	Connected int
	CreatedAt time.Time
}

type Post struct {
	ID        int
	UserID    int
	Username  string
	Title     string
	Body      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
	Comments  []Comment
}

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	PostTitle string
}

// Updated to match notification.go implementation
type Notification struct {
	ID        int64
	UserID    int64
	SenderID  int
	Type      string // Added Type field
	Content   string // Added Content field
	RelatedID int    // Added RelatedID field
	Read      bool   // Changed IsRead to Read
	CreatedAt time.Time
}

type Activity struct {
	Type      string
	Content   string
	Timestamp string
}

type Session struct {
	UserID    int
	Username  string
	CreatedAt time.Time
}

// Struct that will store the content of the json sent by the javascript
type PrivateMessage struct {
	Type     string   `json:"type"`
	Sender   string   `json:"sender"`
	Receiver string   `json:"receiver"`
	Message  string   `json:"message"`
	UserList []string `json:"user_list,omitempty"`
}

type PageData struct {
	Title     string
	Header    string
	Content   interface{}
	IsError   bool
	ErrorCode int
}

type RegisterRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Response struct {
	Username string `json:"username"`
}
