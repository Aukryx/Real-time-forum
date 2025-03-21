package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"user"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PostRequest represents the incoming request structure
type PostRequest struct {
	Title string `json:"title"`
	Body  string `json:"content"`
}

// HandleFetchPosts handles fetching all posts
func HandleFetchPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := db.PostSelectAll()
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// HandleCreatePost handles the creation of new posts
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var postReq PostRequest
	err := json.NewDecoder(r.Body).Decode(&postReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Checking the cookie values
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusUnauthorized)
		return
	}
	userID := db.UserIDWithUUID(cookie.Value)

	// Create a Post from the PostRequest
	post := models.Post{
		UserID: userID,
		Title:  postReq.Title,
		Body:   postReq.Body,
	}

	// Insert the new post into the database
	createdPost, err := db.PostInsert(post.UserID, cookie.Value, post.Title, post.Body)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdPost)
}

// HandleFetchNewPosts handles requests for new posts since a given post ID
func HandleFetchNewPosts(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Get last post ID from query parameter
	lastIDStr := r.URL.Query().Get("lastID")
	lastID, err := strconv.Atoi(lastIDStr)
	if err != nil {
		http.Error(w, "Invalid last ID", http.StatusBadRequest)
		return
	}

	// Query for new posts
	newPosts, err := getNewPosts(lastID)
	if err != nil {
		http.Error(w, "Error fetching new posts", http.StatusInternalServerError)
		return
	}

	// Return posts as JSON
	json.NewEncoder(w).Encode(newPosts)
}

// getNewPosts fetches posts newer than the specified ID
func getNewPosts(lastID int) ([]Post, error) {
	query := `
		SELECT p.id, p.user_id, p.user, p.title, p.body, p.createdAt, p.updatedAt
		FROM post p
		WHERE p.id > ?
		ORDER BY p.id DESC
		LIMIT 20
	`

	rows, err := db.DB.Query(query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Username,
			&post.Title,
			&post.Body,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse datetime strings
		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func HandleFetchPostComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postId"])

	comments, err := db.CommentSelectByPostID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("Error accessing cookie: ", err)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	id := db.UserIDWithUUID(cookie.Value)

	createdComment, err := db.CommentInsert(id, cookie.Value, comment.PostID, comment.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdComment)
}
