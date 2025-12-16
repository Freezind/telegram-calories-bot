package handlers

import (
	"log"

	"github.com/freezind/telegram-calories-bot/internal/services"
	tele "gopkg.in/telebot.v3"
)

// HandleStart handles the /start command.
func HandleStart(c tele.Context) error {
	// Log the received update with user ID
	log.Printf("[UPDATE] UserID=%d Text=%s", c.Sender().ID, c.Text())

	// Send greeting response
	return c.Send(services.Greet())
}

// HandleText handles plain text messages, filtering for "hello".
func HandleText(c tele.Context) error {
	// Log the received update with user ID
	log.Printf("[UPDATE] UserID=%d Text=%s", c.Sender().ID, c.Text())
	// Only respond to exact match "hello" (case-sensitive)
	if c.Text() == "hello" {
		// Send greeting response
		return c.Send(services.Greet())
	}

	// Ignore other messages (no response)
	return nil
}
