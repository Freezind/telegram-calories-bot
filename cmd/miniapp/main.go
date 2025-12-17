package main

import (
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/rs/cors"
	"github.com/freezind/telegram-calories-bot/internal/handlers"
	"github.com/freezind/telegram-calories-bot/internal/middleware"
	"github.com/freezind/telegram-calories-bot/internal/storage"
)

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log completion
		log.Printf("[%s] %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

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

	// API routes for specific log operations (with authentication middleware)
	mux.Handle("/api/logs/", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			logsHandler.UpdateLog(w, r)
		case http.MethodDelete:
			logsHandler.DeleteLog(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Configure CORS for development
	allowedOrigins := []string{"http://localhost:5173"}

	// Add tunnel URL if specified (for cloudflared, ngrok, etc.)
	if tunnelURL := os.Getenv("TUNNEL_URL"); tunnelURL != "" {
		allowedOrigins = append(allowedOrigins, tunnelURL)
		log.Printf("CORS: Added tunnel URL: %s", tunnelURL)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Telegram-Init-Data"},
		AllowCredentials: true,
		Debug:            false, // Disable verbose CORS logging (set true for debugging)
		// Allow requests without Origin header (for Vite proxy and same-origin requests)
		AllowOriginFunc: func(origin string) bool {
			// If no origin header, allow (same-origin or proxied request)
			if origin == "" {
				return true
			}
			// Check against allowed origins
			if slices.Contains(allowedOrigins, origin) {
				return true
			}
			log.Printf("CORS: Blocked origin: %s", origin)
			return false
		},
	})

	// Wrap mux with CORS and logging
	handler := loggingMiddleware(corsHandler.Handler(mux))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("========================================")
	log.Printf("Mini App Backend Server")
	log.Printf("========================================")
	log.Printf("Port: %s", port)
	log.Printf("CORS Origins: %v", allowedOrigins)
	log.Printf("========================================")
	log.Printf("Server ready at http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("API endpoint: http://localhost:%s/api/logs", port)
	log.Printf("========================================")
	log.Printf("")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
