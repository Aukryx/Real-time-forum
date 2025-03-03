package main

import (
	"handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	// Load environment variables
	// lib.LoadEnv(".env")

	// Configure router and server
	mux := setupMux()
	server := setupServer(mux)

	// Settup database
	// db := db.SetupDatabase()

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
	// mux.HandleFunc("/about", handlers.AboutHandler)

	return mux
}

// setupServer configures the HTTP server
func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":8080",
		Handler:           handlers.WithErrorHandling(handler),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
