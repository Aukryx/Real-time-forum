package main

import (
	"config"
	"handlers"
	"log"
	"net/http"
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

	// API routes
	mux.HandleFunc("/api/users", handlers.UserSelectAllHandler)
	mux.HandleFunc("/api/user", handlers.GetUserByIdHandler)
	mux.HandleFunc("/api/posts", handlers.FetchPostsHandler)
	mux.HandleFunc("/api/post", handlers.CreatePostHandler)
	mux.HandleFunc("/api/comments", handlers.FetchPostCommentsHandler)
	mux.HandleFunc("/api/comment", handlers.CreateCommentHandler)

	return mux
}

// setupServer configures the HTTP server
func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":8081",
		Handler:           handlers.WithErrorHandling(handler),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
