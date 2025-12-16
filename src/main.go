// Package main is the entry point for the Telegram calorie estimation bot
package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/freezind/telegram-calories-bot/src/handlers"
	"github.com/freezind/telegram-calories-bot/src/services"
	telebot "gopkg.in/telebot.v3"
)

func main() {
	// Environment variable validation (T017) - fail fast per contracts
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	// Log startup (without exposing secrets)
	log.Println("Starting Calorie Estimation Bot...")
	log.Println("Environment variables validated")

	// Initialize services
	sessionManager := services.NewSessionManager()
	geminiClient, err := services.NewGeminiClient()
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}

	// Start session cleanup goroutine (T018)
	sessionManager.StartCleanupRoutine()
	log.Println("Session cleanup routine started (runs every 5 minutes)")

	// Initialize telebot with settings
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Bot initialized: @%s", bot.Me.Username)
	log.Println("Polling for updates...")

	// Initialize handlers (T024-T033)
	estimateHandler := handlers.NewEstimateHandler(sessionManager, geminiClient)

	// Register command handlers
	bot.Handle("/start", estimateHandler.HandleStart)
	bot.Handle("/estimate", estimateHandler.HandleEstimate)
	bot.Handle(telebot.OnPhoto, estimateHandler.HandlePhoto)
	bot.Handle(telebot.OnDocument, estimateHandler.HandleDocument)

	// Register callback handlers for inline buttons
	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		callbackData := strings.TrimSpace(c.Callback().Data) // Trim whitespace/newlines
		userID := c.Sender().ID

		// Log callback received with details
		log.Printf("[CALLBACK] User %d clicked button. Callback data: '%s' (after trim)", userID, callbackData)

		var err error
		switch callbackData {
		case "re_estimate":
			log.Printf("[CALLBACK] Handling re_estimate for user %d", userID)
			err = estimateHandler.HandleReEstimate(c)
			if err != nil {
				log.Printf("[CALLBACK ERROR] HandleReEstimate failed for user %d: %v", userID, err)
			}
			return err
		case "cancel":
			log.Printf("[CALLBACK] Handling cancel for user %d", userID)
			err = estimateHandler.HandleCancel(c)
			if err != nil {
				log.Printf("[CALLBACK ERROR] HandleCancel failed for user %d: %v", userID, err)
			}
			return err
		default:
			log.Printf("[CALLBACK WARNING] Unknown callback data '%s' from user %d", callbackData, userID)
			return c.Respond(&telebot.CallbackResponse{Text: "Unknown action"})
		}
	})

	log.Println("Handlers registered: /estimate, photo upload, inline buttons")

	// Start bot polling
	bot.Start()
}
