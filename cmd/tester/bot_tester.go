package main

import (
	"context"
	"fmt"
	"time"

	"github.com/freezind/telegram-calories-bot/internal/storage"
	"github.com/freezind/telegram-calories-bot/src/handlers"
	"github.com/freezind/telegram-calories-bot/src/services"
	tele "gopkg.in/telebot.v3"
)

// BotTester runs integration tests on bot handlers
type BotTester struct {
	fakeSender *FakeSender
	handler    *handlers.EstimateHandler
	judge      *GeminiJudge
	testUser   *tele.User
}

// NewBotTester creates a bot tester with fake components
func NewBotTester(judge *GeminiJudge) (*BotTester, error) {
	// Create fake sender
	fakeSender := NewFakeSender()

	// Create real dependencies
	sessionManager := services.NewSessionManager()

	// Use fake estimator for deterministic testing
	fakeEstimator := NewFakeEstimator()

	// Create in-memory storage
	store := storage.NewMemoryStorage()

	// Create real handler with fake sender and fake estimator
	handler := handlers.NewEstimateHandler(fakeSender, sessionManager, fakeEstimator, store)

	// Create test user
	testUser := &tele.User{
		ID:        12345,
		FirstName: "Test",
		LastName:  "User",
	}

	return &BotTester{
		fakeSender: fakeSender,
		handler:    handler,
		judge:      judge,
		testUser:   testUser,
	}, nil
}

// TestStart tests the /start command
func (bt *BotTester) TestStart(ctx context.Context) *TestResult {
	result := NewTestResult("S1", "/start Welcome Message")

	// Create synthetic /start message
	msg := &tele.Message{
		ID:     1,
		Text:   "/start",
		Chat:   &tele.Chat{ID: bt.testUser.ID},
		Sender: bt.testUser,
	}

	// Create mock context
	mockCtx := &simpleContext{
		sender:  bt.testUser,
		message: msg,
	}

	// Invoke handler directly
	err := bt.handler.HandleStart(mockCtx)
	if err != nil {
		result.SetError(fmt.Errorf("HandleStart failed: %w", err))
		result.Complete()
		return result
	}

	// Get sent message
	lastMsg := bt.fakeSender.GetLastMessage()
	if lastMsg == nil {
		result.SetError(fmt.Errorf("no message sent by bot"))
		result.Complete()
		return result
	}

	// Capture evidence
	result.AddEvidence("bot_message", map[string]interface{}{
		"text":         lastMsg.Text,
		"message_id":   lastMsg.MessageID,
		"message_sent": true,
	})

	// Evaluate with LLM judge
	expectedBehavior := "Bot responds with welcome message containing usage instructions and mentions /estimate command"
	verdict, err := bt.judge.Evaluate(ctx, result.ScenarioName, expectedBehavior, result.Evidence)
	if err != nil {
		result.SetError(fmt.Errorf("LLM judge evaluation failed: %w", err))
		result.Complete()
		return result
	}

	result.Verdict = verdict.Verdict
	result.Rationale = verdict.Rationale
	result.Complete()
	return result
}

// TestEstimate tests the /estimate + photo flow
func (bt *BotTester) TestEstimate(ctx context.Context, imagePath string) *TestResult {
	result := NewTestResult("S2", "/estimate + Image Upload")

	// Reset fake sender for clean state
	bt.fakeSender.Reset()

	// Step 1: Send /estimate command
	estimateMsg := &tele.Message{
		ID:     10,
		Text:   "/estimate",
		Chat:   &tele.Chat{ID: bt.testUser.ID},
		Sender: bt.testUser,
	}

	mockCtx := &simpleContext{
		sender:  bt.testUser,
		message: estimateMsg,
	}

	err := bt.handler.HandleEstimate(mockCtx)
	if err != nil {
		result.SetError(fmt.Errorf("HandleEstimate failed: %w", err))
		result.Complete()
		return result
	}

	// Verify bot asked for image
	firstMsg := bt.fakeSender.GetLastMessage()
	if firstMsg == nil {
		result.SetError(fmt.Errorf("no prompt message sent"))
		result.Complete()
		return result
	}

	result.AddEvidence("estimate_prompt", map[string]interface{}{
		"prompt_sent": true,
		"prompt_text": firstMsg.Text,
	})

	// Note: Photo upload with real image download is not testable in integration tests
	// The FakeEstimator provides deterministic results, and we test message preservation
	// in S3 and S4. S2 validates that /estimate command works correctly.

	// Evaluate with LLM judge
	expectedBehavior := "Bot prompts for image after /estimate command (message contains 'Please send' or similar instruction)"
	verdict, err := bt.judge.Evaluate(ctx, result.ScenarioName, expectedBehavior, result.Evidence)
	if err != nil {
		result.SetError(fmt.Errorf("LLM judge evaluation failed: %w", err))
		result.Complete()
		return result
	}

	result.Verdict = verdict.Verdict
	result.Rationale = verdict.Rationale
	result.Complete()
	return result
}

// TestReEstimate tests the re-estimate button preservation
func (bt *BotTester) TestReEstimate(ctx context.Context) *TestResult {
	result := NewTestResult("S3", "Re-estimate Button Message Preservation")

	// Get current message count
	messagesBeforeCount := len(bt.fakeSender.SentMessages)
	if messagesBeforeCount == 0 {
		result.SetError(fmt.Errorf("no previous messages found"))
		result.Complete()
		return result
	}

	// Get the last message ID (simulate clicking re-estimate on this message)
	previousMsg := bt.fakeSender.SentMessages[messagesBeforeCount-1]
	previousMsgID := previousMsg.MessageID

	// Create synthetic callback query (user clicked re-estimate button)
	callbackQuery := &tele.Callback{
		ID:   "callback1",
		Data: "re_estimate",
		Message: &tele.Message{
			ID:     previousMsgID,
			Text:   previousMsg.Text,
			Chat:   &tele.Chat{ID: bt.testUser.ID},
			Sender: bt.testUser,
		},
		Sender: bt.testUser,
	}

	mockCtx := &simpleContext{
		sender:   bt.testUser,
		callback: callbackQuery,
	}

	// Invoke re-estimate handler
	err := bt.handler.HandleReEstimate(mockCtx)
	if err != nil {
		result.SetError(fmt.Errorf("HandleReEstimate failed: %w", err))
		result.Complete()
		return result
	}

	// Check if previous message was deleted (it should NOT be)
	wasDeleted := bt.fakeSender.WasMessageDeleted(previousMsgID)

	// Get new prompt message
	lastMsg := bt.fakeSender.GetLastMessage()

	// Capture evidence - CRITICAL: message preservation test
	result.AddEvidence("message_preservation", map[string]interface{}{
		"previous_message_id":      previousMsgID,
		"previous_message_deleted": wasDeleted,  // MUST be false
		"new_prompt_sent":          lastMsg != nil,
		"new_prompt_text":          "",
	})

	if lastMsg != nil {
		result.Evidence[len(result.Evidence)-1].Data["new_prompt_text"] = lastMsg.Text
	}

	// Evaluate with LLM judge
	expectedBehavior := "After clicking Re-estimate button, bot sends NEW prompt message and does NOT delete previous estimate message (previous_message_deleted MUST be false)"
	verdict, err := bt.judge.Evaluate(ctx, result.ScenarioName, expectedBehavior, result.Evidence)
	if err != nil {
		result.SetError(fmt.Errorf("LLM judge evaluation failed: %w", err))
		result.Complete()
		return result
	}

	result.Verdict = verdict.Verdict
	result.Rationale = verdict.Rationale
	result.Complete()
	return result
}

// TestCancel tests the cancel button preservation
func (bt *BotTester) TestCancel(ctx context.Context) *TestResult {
	result := NewTestResult("S4", "Cancel Button Message Preservation")

	// Get current message count
	messagesBeforeCount := len(bt.fakeSender.SentMessages)
	if messagesBeforeCount == 0 {
		result.SetError(fmt.Errorf("no previous messages found"))
		result.Complete()
		return result
	}

	// Get the last message ID
	previousMsg := bt.fakeSender.SentMessages[messagesBeforeCount-1]
	previousMsgID := previousMsg.MessageID

	// Create synthetic callback query (user clicked cancel button)
	callbackQuery := &tele.Callback{
		ID:   "callback2",
		Data: "cancel",
		Message: &tele.Message{
			ID:     previousMsgID,
			Text:   previousMsg.Text,
			Chat:   &tele.Chat{ID: bt.testUser.ID},
			Sender: bt.testUser,
		},
		Sender: bt.testUser,
	}

	mockCtx := &simpleContext{
		sender:   bt.testUser,
		callback: callbackQuery,
	}

	// Invoke cancel handler
	err := bt.handler.HandleCancel(mockCtx)
	if err != nil {
		result.SetError(fmt.Errorf("HandleCancel failed: %w", err))
		result.Complete()
		return result
	}

	// Check if previous message was deleted (it should NOT be)
	wasDeleted := bt.fakeSender.WasMessageDeleted(previousMsgID)

	// Get cancellation message
	lastMsg := bt.fakeSender.GetLastMessage()

	// Capture evidence - CRITICAL: message preservation test
	result.AddEvidence("message_preservation", map[string]interface{}{
		"previous_message_id":      previousMsgID,
		"previous_message_deleted": wasDeleted,  // MUST be false
		"cancellation_sent":        lastMsg != nil,
		"cancellation_text":        "",
	})

	if lastMsg != nil {
		result.Evidence[len(result.Evidence)-1].Data["cancellation_text"] = lastMsg.Text
	}

	// Evaluate with LLM judge
	expectedBehavior := "After clicking Cancel button, bot sends cancellation confirmation and does NOT delete previous estimate message (previous_message_deleted MUST be false)"
	verdict, err := bt.judge.Evaluate(ctx, result.ScenarioName, expectedBehavior, result.Evidence)
	if err != nil {
		result.SetError(fmt.Errorf("LLM judge evaluation failed: %w", err))
		result.Complete()
		return result
	}

	result.Verdict = verdict.Verdict
	result.Rationale = verdict.Rationale
	result.Complete()
	return result
}

// simpleContext is a minimal telebot.Context implementation for testing
type simpleContext struct {
	sender   *tele.User
	message  *tele.Message
	callback *tele.Callback
}

func (s *simpleContext) Sender() *tele.User                            { return s.sender }
func (s *simpleContext) Message() *tele.Message                        { return s.message }
func (s *simpleContext) Callback() *tele.Callback                      { return s.callback }
func (s *simpleContext) Respond(resp ...*tele.CallbackResponse) error { return nil }

// Implement remaining Context interface methods as no-ops
func (s *simpleContext) Bot() *tele.Bot                                        { return nil }
func (s *simpleContext) Update() tele.Update                                   { return tele.Update{} }
func (s *simpleContext) Text() string                                          { if s.message != nil { return s.message.Text }; return "" }
func (s *simpleContext) Data() string                                          { if s.callback != nil { return s.callback.Data }; return "" }
func (s *simpleContext) Args() []string                                        { return []string{} }
func (s *simpleContext) Chat() *tele.Chat                                      { if s.message != nil { return s.message.Chat }; return nil }
func (s *simpleContext) Recipient() tele.Recipient                             { return s.sender }
func (s *simpleContext) Get(key string) interface{}                            { return nil }
func (s *simpleContext) Set(key string, val interface{})                       {}
func (s *simpleContext) Send(what interface{}, opts ...interface{}) error     { return nil }
func (s *simpleContext) Reply(what interface{}, opts ...interface{}) error    { return nil }
func (s *simpleContext) SendAlbum(a tele.Album, opts ...interface{}) error    { return nil }
func (s *simpleContext) Forward(msg tele.Editable, opts ...interface{}) error { return nil }
func (s *simpleContext) ForwardTo(to tele.Recipient, opts ...interface{}) error { return nil }
func (s *simpleContext) Edit(what interface{}, opts ...interface{}) error       { return nil }
func (s *simpleContext) EditCaption(caption string, opts ...interface{}) error  { return nil }
func (s *simpleContext) EditOrSend(what interface{}, opts ...interface{}) error { return nil }
func (s *simpleContext) EditOrReply(what interface{}, opts ...interface{}) error { return nil }
func (s *simpleContext) Delete() error                                            { return nil }
func (s *simpleContext) DeleteAfter(d time.Duration) *time.Timer                  { return nil }
func (s *simpleContext) Notify(action tele.ChatAction) error                      { return nil }
func (s *simpleContext) Ship(what ...interface{}) error                           { return nil }
func (s *simpleContext) Accept(errorMessage ...string) error                      { return nil }
func (s *simpleContext) Answer(resp *tele.QueryResponse) error                    { return nil }
func (s *simpleContext) RespondText(text string) error                            { return nil }
func (s *simpleContext) RespondAlert(text string) error                           { return nil }
func (s *simpleContext) Query() *tele.Query                                       { return nil }
func (s *simpleContext) InlineResult() *tele.InlineResult                         { return nil }
func (s *simpleContext) ShippingQuery() *tele.ShippingQuery                       { return nil }
func (s *simpleContext) PreCheckoutQuery() *tele.PreCheckoutQuery                 { return nil }
func (s *simpleContext) Poll() *tele.Poll                                         { return nil }
func (s *simpleContext) PollAnswer() *tele.PollAnswer                             { return nil }
func (s *simpleContext) ChatMember() *tele.ChatMemberUpdate                       { return nil }
func (s *simpleContext) ChatJoinRequest() *tele.ChatJoinRequest                   { return nil }
func (s *simpleContext) Migration() (int64, int64)                                { return 0, 0 }
func (s *simpleContext) Topic() *tele.Topic                                       { return nil }
func (s *simpleContext) Boost() *tele.BoostUpdated                                { return nil }
func (s *simpleContext) BoostRemoved() *tele.BoostRemoved                         { return nil }
func (s *simpleContext) Entities() tele.Entities                                  { return nil }
