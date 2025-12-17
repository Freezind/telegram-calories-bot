package bot

import (
	tele "gopkg.in/telebot.v3"
)

// Sender is a minimal interface for sending messages, editing, and deleting
// This abstraction allows testing without a real Telegram bot
type Sender interface {
	// Send sends a message to the recipient
	Send(to tele.Recipient, what interface{}, opts ...interface{}) (*tele.Message, error)

	// Delete removes a message
	Delete(msg tele.Editable) error

	// Respond answers a callback query
	Respond(callback *tele.Callback, resp ...*tele.CallbackResponse) error

	// FileByID retrieves file information by file ID
	FileByID(fileID string) (tele.File, error)

	// GetFileURL returns the download URL for a file
	GetFileURL(file tele.File) string
}

// TelebotSender wraps a real telebot.Bot to implement the Sender interface
type TelebotSender struct {
	bot *tele.Bot
}

// NewTelebotSender creates a production Sender wrapping a real bot
func NewTelebotSender(bot *tele.Bot) Sender {
	return &TelebotSender{bot: bot}
}

func (s *TelebotSender) Send(to tele.Recipient, what interface{}, opts ...interface{}) (*tele.Message, error) {
	return s.bot.Send(to, what, opts...)
}

func (s *TelebotSender) Delete(msg tele.Editable) error {
	return s.bot.Delete(msg)
}

func (s *TelebotSender) Respond(callback *tele.Callback, resp ...*tele.CallbackResponse) error {
	return s.bot.Respond(callback, resp...)
}

func (s *TelebotSender) FileByID(fileID string) (tele.File, error) {
	return s.bot.FileByID(fileID)
}

func (s *TelebotSender) GetFileURL(file tele.File) string {
	return s.bot.URL + "/file/bot" + s.bot.Token + "/" + file.FilePath
}
