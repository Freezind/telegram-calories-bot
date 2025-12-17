package main

import (
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/rs/cors"
	tele "gopkg.in/telebot.v3"

	apihandlers "github.com/freezind/telegram-calories-bot/internal/handlers"
	"github.com/freezind/telegram-calories-bot/internal/middleware"
	"github.com/freezind/telegram-calories-bot/internal/storage"
	"github.com/freezind/telegram-calories-bot/src/bot"
	bothandlers "github.com/freezind/telegram-calories-bot/src/handlers"
	"github.com/freezind/telegram-calories-bot/src/services"
)

// Logging middleware for HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[HTTP %s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("[HTTP %s] %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	log.Println("========================================")
	log.Println("Unified Telegram Calorie Bot + Mini App")
	log.Println("========================================")

	// ====================================
	// 1. Initialize shared storage
	// ====================================
	store := storage.NewMemoryStorage()
	log.Println("[STORAGE] ✓ Shared MemoryStorage initialized")

	// ====================================
	// 2. Initialize Telegram Bot (Spec 002)
	// ====================================
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// GEMINI_API_KEY is read by NewGeminiClient() from environment
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("❌ GEMINI_API_KEY environment variable is required")
	}

	// Create bot instance
	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	tgBot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("❌ Failed to create bot: %v", err)
	}

	// Wrap bot as Sender
	sender := bot.NewTelebotSender(tgBot)

	// Initialize bot dependencies
	sessionManager := services.NewSessionManager()
	geminiClient, err := services.NewGeminiClient()
	if err != nil {
		log.Fatalf("❌ Failed to initialize Gemini client: %v", err)
	}
	estimator := services.NewGeminiEstimator(geminiClient)
	estimateHandler := bothandlers.NewEstimateHandler(sender, sessionManager, estimator, store)

	// Register bot command handlers
	tgBot.Handle("/start", estimateHandler.HandleStart)
	tgBot.Handle("/estimate", estimateHandler.HandleEstimate)
	tgBot.Handle(tele.OnPhoto, estimateHandler.HandlePhoto)
	tgBot.Handle(tele.OnDocument, estimateHandler.HandleDocument)

	// Register callback handlers for inline buttons
	tgBot.Handle(tele.OnCallback, func(c tele.Context) error {
		callbackData := strings.TrimSpace(c.Callback().Data)
		userID := c.Sender().ID

		log.Printf("[BOT CALLBACK] User %d clicked button: '%s'", userID, callbackData)

		switch callbackData {
		case "re_estimate":
			return estimateHandler.HandleReEstimate(c)
		case "cancel":
			return estimateHandler.HandleCancel(c)
		default:
			log.Printf("[BOT CALLBACK] Unknown callback: '%s'", callbackData)
			return c.Respond(&tele.CallbackResponse{Text: "Unknown action"})
		}
	})

	// Start session cleanup routine
	sessionManager.StartCleanupRoutine()
	log.Println("[BOT] ✓ Telegram bot initialized (session cleanup routine started)")

	// ====================================
	// 3. Initialize HTTP API Server (Spec 003)
	// ====================================
	logsHandler := apihandlers.NewLogsHandler(store)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes with authentication middleware
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

	// Configure CORS
	allowedOrigins := []string{
		"http://localhost:5173",
		"https://telegram-calories-bot.pages.dev",
	}
	if tunnelURL := os.Getenv("TUNNEL_URL"); tunnelURL != "" {
		allowedOrigins = append(allowedOrigins, tunnelURL)
		log.Printf("[HTTP] CORS: Added tunnel URL: %s", tunnelURL)
	}
	if miniappURL := os.Getenv("MINIAPP_URL"); miniappURL != "" {
		allowedOrigins = append(allowedOrigins, miniappURL)
		log.Printf("[HTTP] CORS: Added miniapp URL: %s", miniappURL)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Telegram-Init-Data"},
		AllowCredentials: true,
		Debug:            false,
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return true
			}
			if slices.Contains(allowedOrigins, origin) {
				return true
			}
			log.Printf("[HTTP] CORS: Blocked origin: %s", origin)
			return false
		},
	})

	handler := loggingMiddleware(corsHandler.Handler(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("[HTTP] ✓ HTTP API server initialized")
	log.Printf("[HTTP] Port: %s", port)
	log.Printf("[HTTP] CORS Origins: %v", allowedOrigins)

	// ====================================
	// 4. Start both services concurrently
	// ====================================
	log.Println("========================================")
	log.Println("Starting services...")
	log.Println("========================================")

	// Start HTTP server in goroutine
	go func() {
		log.Printf("[HTTP] 🚀 HTTP server listening on http://localhost:%s", port)
		log.Printf("[HTTP] Health check: http://localhost:%s/health", port)
		log.Printf("[HTTP] API endpoint: http://localhost:%s/api/logs", port)
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			log.Fatalf("❌ HTTP server failed: %v", err)
		}
	}()

	// Start Telegram bot (blocking)
	log.Println("[BOT] 🚀 Telegram bot started")
	tgBot.Start()
}
