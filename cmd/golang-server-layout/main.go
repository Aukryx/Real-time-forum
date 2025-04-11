package main

import (
	"config"
	"encoding/json"
	"handlers"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Initialize configuration
	log.Println("Application starting...")
	config.Initialize()
	log.Println("Configuration loaded.")

	// Configure router and server
	mux := setupMux()
	server := setupServer(mux)

	log.Printf("Server starting on http://localhost%s...", server.Addr)

	// Start HTTP server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server stopped.")
}

// setupMux configures routes
func setupMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files with better path handling
	staticPath := "../../web/static" // Original path
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		// If running from root directory (as on Render)
		staticPath = "web/static"
		log.Printf("Using alternate static path: %s", staticPath)
	}

	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	mux.HandleFunc("/", handlers.IndexHandler)

	// Authentication routes
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/api/check-session", handlers.CheckSession)

	// Replace the WebSocket route with a conditional
	if os.Getenv("PORT") != "" {
		// On Render - return a simple JSON response
		mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // Return 200 instead of 400
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "disabled",
				"message": "WebSockets are disabled in demo mode",
			})
		})
	} else {
		// Local development - use real WebSockets
		mux.HandleFunc("/ws", handlers.HandleConnection)
	}

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
	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8060"
	}

	return &http.Server{
		Addr:              ":" + port,
		Handler:           handlers.WithErrorHandling(handler),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
