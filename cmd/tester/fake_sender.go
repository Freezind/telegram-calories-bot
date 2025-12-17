package main

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

// FakeSender records all outgoing actions for testing
type FakeSender struct {
	SentMessages     []SentMessage
	DeletedMessages  []int          // Message IDs that were deleted
	EditedMessages   []SentMessage  // Messages that were edited
	CallbackResponses []string       // Callback responses sent
}

// SentMessage captures a message sent by the bot
type SentMessage struct {
	ChatID      int64
	Text        string
	MessageID   int
	ReplyMarkup *tele.ReplyMarkup
}

// NewFakeSender creates a new fake sender for testing
func NewFakeSender() *FakeSender {
	return &FakeSender{
		SentMessages:      []SentMessage{},
		DeletedMessages:   []int{},
		EditedMessages:    []SentMessage{},
		CallbackResponses: []string{},
	}
}

// Send captures sent messages
func (f *FakeSender) Send(to tele.Recipient, what interface{}, opts ...interface{}) (*tele.Message, error) {
	// Extract chat ID
	var chatID int64
	if user, ok := to.(*tele.User); ok {
		chatID = user.ID
	} else if chat, ok := to.(*tele.Chat); ok {
		chatID = chat.ID
	}

	// Extract text
	text := ""
	switch v := what.(type) {
	case string:
		text = v
	default:
		text = fmt.Sprintf("%v", what)
	}

	// Extract reply markup
	var replyMarkup *tele.ReplyMarkup
	for _, opt := range opts {
		if rm, ok := opt.(*tele.ReplyMarkup); ok {
			replyMarkup = rm
			break
		}
	}

	// Create message ID
	msgID := len(f.SentMessages) + 1

	// Record sent message
	sentMsg := SentMessage{
		ChatID:      chatID,
		Text:        text,
		MessageID:   msgID,
		ReplyMarkup: replyMarkup,
	}
	f.SentMessages = append(f.SentMessages, sentMsg)

	// Return mock message
	return &tele.Message{
		ID:   msgID,
		Text: text,
		Chat: &tele.Chat{ID: chatID},
	}, nil
}

// Delete records message deletions
func (f *FakeSender) Delete(msg tele.Editable) error {
	// Extract message ID
	if m, ok := msg.(*tele.Message); ok {
		f.DeletedMessages = append(f.DeletedMessages, m.ID)
	}
	return nil
}

// Respond records callback responses
func (f *FakeSender) Respond(callback *tele.Callback, resp ...*tele.CallbackResponse) error {
	if len(resp) > 0 && resp[0] != nil {
		f.CallbackResponses = append(f.CallbackResponses, resp[0].Text)
	}
	return nil
}

// FileByID returns mock file info
func (f *FakeSender) FileByID(fileID string) (tele.File, error) {
	return tele.File{
		FileID:   fileID,
		FilePath: "test/path/" + fileID,
	}, nil
}

// GetFileURL returns a mock file URL
func (f *FakeSender) GetFileURL(file tele.File) string {
	return "https://api.telegram.org/file/botTOKEN/" + file.FilePath
}

// GetLastMessage returns the most recently sent message
func (f *FakeSender) GetLastMessage() *SentMessage {
	if len(f.SentMessages) == 0 {
		return nil
	}
	return &f.SentMessages[len(f.SentMessages)-1]
}

// WasMessageDeleted checks if a message ID was deleted
func (f *FakeSender) WasMessageDeleted(msgID int) bool {
	for _, id := range f.DeletedMessages {
		if id == msgID {
			return true
		}
	}
	return false
}

// Reset clears all recorded actions
func (f *FakeSender) Reset() {
	f.SentMessages = []SentMessage{}
	f.DeletedMessages = []int{}
	f.EditedMessages = []SentMessage{}
	f.CallbackResponses = []string{}
}
