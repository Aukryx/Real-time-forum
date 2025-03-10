package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"log"
	"models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ImagePath *string   `json:"image"`
}

// PostRequest represents the incoming request structure
type PostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// HandleFetchPosts handles fetching all posts
func HandleFetchPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := db.PostSelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// HandleCreatePost handles the creation of new posts
// Then modify HandleCreatePost to use this structure:
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var postReq PostRequest
	err := json.NewDecoder(r.Body).Decode(&postReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a Post from the PostRequest
	post := models.Post{
		Title: postReq.Title,
		Body:  postReq.Body,
	}

	fmt.Println("post: ", post)

	// Checking the cookie values
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("Error accessing cookie: ", err)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	id := db.UserIDWithUUID(cookie.Value)

	// Ensure you're using the correct column names
	createdPost, err := db.PostInsert(id, post.Title, post.Body, post.ImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	lastIDStr := r.URL.Query().Get("lastId")
	lastID, err := strconv.Atoi(lastIDStr)
	if err != nil {
		// Default to 0 if no valid ID is provided
		lastID = 0
	}

	// Query for new posts
	posts, err := getNewPosts(lastID)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	// Return posts as JSON
	response := map[string]interface{}{
		"posts": posts,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// getNewPosts fetches posts newer than the specified ID
func getNewPosts(lastID int) ([]Post, error) {
	query := `
		SELECT id, user_id, title, body, createdAt, updatedAt, image
		FROM post
		WHERE id > ?
		ORDER BY id DESC
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
			&post.Title,
			&post.Body,
			&createdAtStr,
			&updatedAtStr,
			&post.ImagePath,
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
	var comment db.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdComment, err := db.CommentInsert(comment.UserID, comment.PostID, comment.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdComment)
}
