// internal/handlers/api.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"db"
	"models"
)

func UserSelectAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := db.UserSelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from query parameter
	idStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := db.UserSelectByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func FetchPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := db.PostSelectAll()
	if err != nil {
		// Add more detailed logging
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Failed to fetch posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify posts are not empty
	if len(posts) == 0 {
		log.Println("No posts found")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add logging to debug
	log.Printf("Received post: %+v", post)

	// Ensure you're passing the correct parameters
	createdPost, err := db.PostInsert(post.UserID, "", post.Title, post.Body, post.ImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdPost)
}

// internal/handlers/api.go
// Fix the FetchPostCommentsHandler function to correctly extract the post ID

func FetchPostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract post ID from URL path - format: /api/posts/{postId}/comments
	path := r.URL.Path

	// Split the path and validate it has enough parts
	pathParts := strings.Split(path, "/")

	// Validate path has enough parts (should be like ["", "api", "posts", "{postId}", "comments"])
	if len(pathParts) < 5 || pathParts[1] != "api" || pathParts[2] != "posts" || pathParts[4] != "comments" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// The postId should be at index 3 in the path parts array
	postIDStr := pathParts[3]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Now fetch comments with the extracted postID
	comments, err := db.CommentSelectByPostID(postID)
	if err != nil {
		http.Error(w, "Error fetching comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract post ID from URL path - format: /api/posts/{postId}/comments
	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	// Validate path structure
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// The postId should be at index 3 in the path parts array
	postIDStr := pathParts[3]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the postID from the URL
	comment.PostID = postID

	// Checking the cookie values
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusUnauthorized)
		return
	}
	userID := db.UserIDWithUUID(cookie.Value)

	// Insert the new comment into the database
	createdComment, err := db.CommentInsert(userID, cookie.Value, comment.PostID, comment.Body)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdComment)
}
