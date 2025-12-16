package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/freezind/telegram-calories-bot/internal/handlers"
	"github.com/freezind/telegram-calories-bot/internal/middleware"
	"github.com/freezind/telegram-calories-bot/internal/storage"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize handlers
	logsHandler := handlers.NewLogsHandler(store)

	// Create HTTP router
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes (with authentication middleware)
	mux.Handle("/api/logs", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			logsHandler.ListLogs(w, r)
		case http.MethodPost:
			logsHandler.CreateLog(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Configure CORS for localhost:5173 (Vite dev server)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Telegram-Init-Data"},
		AllowCredentials: true,
	})

	// Wrap mux with CORS
	handler := corsHandler.Handler(mux)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Mini App backend server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
