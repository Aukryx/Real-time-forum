package main

import (
	"config"
	"handlers"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Initialize configuration
	config.Initialize()

	// Configure router and server
	mux := setupMux()
	server := setupServer(mux)

	log.Printf("Server starting on http://localhost%s...", server.Addr)

	// Start HTTP server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// setupMux configures routes
func setupMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("../../web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	mux.HandleFunc("/", handlers.IndexHandler)

	// Authentication routes
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/api/check-session", handlers.CheckSession)
	mux.HandleFunc("/ws", handlers.HandleConnection)

	// API routes
	mux.HandleFunc("/api/users", handlers.GetConnectedAndDisconnectedUsers)
	mux.HandleFunc("/api/user", handlers.GetUserByIdHandler)
	mux.HandleFunc("/api/post", handlers.CreatePostHandler)
	mux.HandleFunc("/api/posts", handlers.HandleFetchPosts)
	mux.HandleFunc("/api/postCreation", handlers.HandleCreatePost)
	mux.HandleFunc("/api/posts/new", handlers.HandleFetchNewPosts)
	mux.HandleFunc("/api/navbar", handlers.NavbarHandler)

	// Handle comment-related routes
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Check if this matches the comments pattern: /api/posts/{postId}/comments
		if strings.Contains(path, "/comments") {
			if r.Method == http.MethodGet {
				handlers.FetchPostCommentsHandler(w, r)
				return
			} else if r.Method == http.MethodPost {
				handlers.CreateCommentHandler(w, r)
				return
			}
		}

		// If we get here, it wasn't a comments request
		http.NotFound(w, r)
	})

	// Session management
	mux.HandleFunc("/logout", handlers.LogOutHandler)

	return mux
}

// setupServer configures the HTTP server
func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":8050",
		Handler:           handlers.WithErrorHandling(handler),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
