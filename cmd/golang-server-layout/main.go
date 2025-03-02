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
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("/static/uploads/", http.StripPrefix("/static/uploads/", http.FileServer(http.Dir("web/static/uploads/"))))

	// Define routes
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/ws/register", handlers.RegisterHandler)
	// mux.HandleFunc("/about", handlers.AboutHandler)

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
