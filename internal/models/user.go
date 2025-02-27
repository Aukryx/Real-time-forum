package models

import (
	"time"
)

type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	NickName  string
	Password  string
	Role      string
	CreatedAt time.Time
}

type Post struct {
	ID            int
	UserID        int
	Title         string
	Body          string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	User          User
	ImagePath     string
	Categories    []Category
	Comments      []Comment
	LikeDislike   []LikeDislike
	LikeCount     int
	DislikesCount int
}

type Comment struct {
	ID            int
	PostID        int
	UserID        int
	Content       string
	CreatedAt     string
	UpdatedAt     string
	LikeDislike   []LikeDislike
	Username      string
	LikesCount    int
	DislikesCount int
	PostTitle     string
}

type Category struct {
	ID   int
	Name string
}

type PostCategory struct {
	PostID     int
	CategoryId int
}

// Updated to match like_dislike.go implementation
type LikeDislike struct {
	ID        int
	UserID    int
	PostID    int
	CommentID int
	Status    int // Changed from IsLike boolean to Status int
	CreatedAt time.Time
	PostTitle string
	Username  string
}

type Image struct {
	ID        int
	PostID    int
	FilePath  string
	FileSize  int
	CreatedAt time.Time
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

type PrivateMessage struct {
	ID         int
	SenderID   int
	ReceiverID int
	Message    string
	CreatedAt  time.Time
	Read       bool
}
