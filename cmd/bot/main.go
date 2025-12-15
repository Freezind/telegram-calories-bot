package main

import (
	"log"
	"os"
	"time"

	"github.com/freezind/telegram-calories-bot/internal/handlers"
	tele "gopkg.in/telebot.v3"
)

func main() {
	// Load bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Create bot with long polling
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Register /start command handler
	bot.Handle("/start", handlers.HandleStart)

	// Register text message handler
	bot.Handle(tele.OnText, handlers.HandleText)

	log.Println("Bot started successfully")
	bot.Start()
}
