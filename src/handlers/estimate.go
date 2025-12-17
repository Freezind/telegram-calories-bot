// Package handlers implements Telegram bot message and command handlers
package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	internalmodels "github.com/freezind/telegram-calories-bot/internal/models"
	"github.com/freezind/telegram-calories-bot/src/models"
	"github.com/freezind/telegram-calories-bot/src/services"
	telebot "gopkg.in/telebot.v3"
)

// EstimateHandler handles bot commands and interactions
// Handles /start, /estimate, image uploads, and inline buttons
type EstimateHandler struct {
	sessionManager *services.SessionManager
	geminiClient   *services.GeminiClient
	storage        LogStorage // Interface for log persistence (shared with miniapp)
}

// LogStorage defines the interface for storing calorie logs
// This matches internal/storage/interface.go
type LogStorage interface {
	CreateLog(userID int64, log *internalmodels.Log) error
}

// HandleStart handles the /start command (T086)
// Sends welcome message with bot introduction and usage instructions
func (h *EstimateHandler) HandleStart(c telebot.Context) error {
	welcomeMsg := models.FormatWelcomeMessage()
	_, err := c.Bot().Send(c.Sender(), welcomeMsg)
	if err != nil {
		log.Printf("[HANDLER ERROR] Failed to send welcome message: %v", err)
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	log.Printf("[HANDLER] HandleStart sent welcome message to user %d", c.Sender().ID)
	return nil
}

// NewEstimateHandler creates a new EstimateHandler instance
func NewEstimateHandler(sm *services.SessionManager, gc *services.GeminiClient, storage LogStorage) *EstimateHandler {
	return &EstimateHandler{
		sessionManager: sm,
		geminiClient:   gc,
		storage:        storage,
	}
}

// HandleEstimate handles the /estimate command
// Flow: User sends /estimate → Bot prompts for image → State: AwaitingImage
func (h *EstimateHandler) HandleEstimate(c telebot.Context) error {
	userID := c.Sender().ID

	// Update session state to AwaitingImage
	h.sessionManager.UpdateSession(userID, models.StateAwaitingImage)

	// Send prompt message with Cancel button (FR-002, FR-007)
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Cancel", "cancel")
	markup.Inline(
		markup.Row(btnCancel),
	)

	msg, err := c.Bot().Send(c.Sender(), "📸 Please send a food image for calorie estimation", markup)
	if err != nil {
		return fmt.Errorf("failed to send prompt: %w", err)
	}

	// Store message ID for later editing/deletion
	h.sessionManager.SetMessageID(userID, msg.ID)

	return nil
}

// HandleDocument handles document uploads (for PNG, WebP original files)
// Telegram compresses photos to JPEG, so original PNG/WebP must be sent as documents
func (h *EstimateHandler) HandleDocument(c telebot.Context) error {
	userID := c.Sender().ID

	// Check session state - only process if AwaitingImage
	session := h.sessionManager.GetSession(userID)
	if session.State != models.StateAwaitingImage {
		return nil
	}

	// Check if document is an image
	doc := c.Message().Document
	if doc == nil {
		return nil
	}

	// Validate image format (T026 - FR-003)
	if !isValidImageFormat(doc.MIME) {
		return h.sendError(c, "Unsupported format. Please send JPEG, PNG, or WebP images only.")
	}

	// Process the document as an image
	return h.processImage(c, doc.FileID, doc.MIME)
}

// HandlePhoto handles photo uploads (T025)
// Photos are always JPEG in Telegram (compressed)
func (h *EstimateHandler) HandlePhoto(c telebot.Context) error {
	userID := c.Sender().ID

	// Check session state - only process if AwaitingImage
	session := h.sessionManager.GetSession(userID)
	if session.State != models.StateAwaitingImage {
		return nil
	}

	photo := c.Message().Photo
	if photo == nil {
		return nil
	}

	// Process as JPEG (Telegram default)
	return h.processImage(c, photo.FileID, "image/jpeg")
}

// processImage handles the common image processing logic for both photos and documents
func (h *EstimateHandler) processImage(c telebot.Context, fileID, mimeType string) error {

	userID := c.Sender().ID
	ctx := context.Background()

	// Update state to Processing
	h.sessionManager.UpdateSession(userID, models.StateProcessing)

	// Send processing message
	processingMsg, err := c.Bot().Send(c.Sender(), "⏳ Analyzing your image...")
	if err != nil {
		log.Printf("Failed to send processing message: %v", err)
	}

	// Download image from Telegram
	file, err := c.Bot().FileByID(fileID)
	if err != nil {
		h.sessionManager.UpdateSession(userID, models.StateIdle)
		return h.sendError(c, "Failed to download image. Please try again.")
	}

	// Fetch file content
	fileURL := c.Bot().URL + "/file/bot" + c.Bot().Token + "/" + file.FilePath
	log.Printf("fileURL: %s", fileURL)
	// #nosec G107 - URL is constructed from trusted Telegram Bot API response
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		h.sessionManager.UpdateSession(userID, models.StateIdle)
		return h.sendError(c, "Failed to create request. Please try again.")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.sessionManager.UpdateSession(userID, models.StateIdle)
		return h.sendError(c, "Failed to fetch image. Please try again.")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		h.sessionManager.UpdateSession(userID, models.StateIdle)
		return h.sendError(c, "Failed to read image. Please try again.")
	}

	// Call Gemini Vision API (T028)
	result, err := h.geminiClient.EstimateCalories(ctx, imageBytes, mimeType)
	if err != nil {
		log.Printf("error when call gemini API: %v", err)
		h.sessionManager.UpdateSession(userID, models.StateIdle)
		// Delete processing message
		if processingMsg != nil {
			if delErr := c.Bot().Delete(processingMsg); delErr != nil {
				log.Printf("Failed to delete processing message: %v", delErr)
			}
		}
		return h.sendError(c, "API error. Please try again later.") // T033
	}

	// Check if food was detected (T031 - FR-014)
	if !result.HasFood() {
		h.sessionManager.UpdateSession(userID, models.StateIdle)
		// Delete processing message
		if processingMsg != nil {
			if delErr := c.Bot().Delete(processingMsg); delErr != nil {
				log.Printf("Failed to delete processing message: %v", delErr)
			}
		}
		return h.sendError(c, "No food detected in image. Please send an image containing food.")
	}

	// Delete processing message
	if processingMsg != nil {
		if delErr := c.Bot().Delete(processingMsg); delErr != nil {
			log.Printf("Failed to delete processing message: %v", delErr)
		}
	}

	// Format and send result (T030 - FR-006)
	formattedResult := models.FormatResult(result)

	// Create inline keyboard with Re-estimate and Cancel buttons (T029 - FR-008, FR-009)
	markup := &telebot.ReplyMarkup{}
	btnReEstimate := markup.Data("Re-estimate", "re_estimate")
	btnCancel := markup.Data("Cancel", "cancel")
	markup.Inline(
		markup.Row(btnReEstimate, btnCancel),
	)

	_, err = c.Bot().Send(c.Sender(), formattedResult, markup)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	// Store the log entry in shared storage (visible in miniapp)
	if h.storage != nil {
		logEntry := &internalmodels.Log{
			FoodItems:  result.FoodItems,
			Calories:   result.Calories,
			Confidence: internalmodels.ConfidenceLevel(result.Confidence),
			Timestamp:  time.Now(),
		}

		if err := h.storage.CreateLog(userID, logEntry); err != nil {
			// Log error but don't fail the user's request
			log.Printf("[HANDLER ERROR] Failed to save log entry for user %d: %v", userID, err)
		} else {
			log.Printf("[HANDLER] ✓ Log entry saved for user %d: %d kcal, %d items", userID, result.Calories, len(result.FoodItems))
		}
	}

	// Keep session in AwaitingImage state for potential Re-estimate
	h.sessionManager.UpdateSession(userID, models.StateAwaitingImage)

	return nil
}

// HandleReEstimate handles the Re-estimate button click (User Story 2)
// T089: Modified to preserve previous message (no deletion)
func (h *EstimateHandler) HandleReEstimate(c telebot.Context) error {
	userID := c.Sender().ID
	log.Printf("[HANDLER] HandleReEstimate called for user %d", userID)

	// Update callback to show feedback
	if err := c.Respond(&telebot.CallbackResponse{Text: "Send another image"}); err != nil {
		log.Printf("[HANDLER ERROR] Failed to respond to re_estimate callback for user %d: %v", userID, err)
	}

	// T089: Keep previous result visible - DO NOT delete message
	// Previous implementation deleted the message with c.Delete()
	// Now we preserve conversation history

	// Update state to AwaitingImage
	h.sessionManager.UpdateSession(userID, models.StateAwaitingImage)

	// Send new prompt
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Cancel", "cancel")
	markup.Inline(
		markup.Row(btnCancel),
	)

	msg, err := c.Bot().Send(c.Sender(), "📸 Please send another food image", markup)
	if err != nil {
		log.Printf("[HANDLER ERROR] Failed to send re-estimate prompt for user %d: %v", userID, err)
		return fmt.Errorf("failed to send re-estimate prompt: %w", err)
	}

	h.sessionManager.SetMessageID(userID, msg.ID)

	log.Printf("[HANDLER] HandleReEstimate completed successfully for user %d", userID)
	return nil
}

// HandleCancel handles the Cancel button click (User Story 3)
// T090: Modified to preserve previous messages (no deletion)
func (h *EstimateHandler) HandleCancel(c telebot.Context) error {
	userID := c.Sender().ID
	log.Printf("[HANDLER] HandleCancel called for user %d", userID)

	// Update callback to show feedback
	if err := c.Respond(&telebot.CallbackResponse{Text: "Estimation canceled"}); err != nil {
		log.Printf("[HANDLER ERROR] Failed to respond to cancel callback for user %d: %v", userID, err)
	}

	// T090: Keep previous messages visible - DO NOT delete message
	// Previous implementation deleted the message with c.Delete()
	// Now we preserve conversation history

	// Clean up session
	h.sessionManager.DeleteSession(userID)

	// Send cancellation confirmation (FR-013)
	_, err := c.Bot().Send(c.Sender(), "Estimation canceled. Use /estimate to start again.")
	if err != nil {
		log.Printf("[HANDLER ERROR] Failed to send cancellation message for user %d: %v", userID, err)
		return fmt.Errorf("failed to send cancellation message: %w", err)
	}

	log.Printf("[HANDLER] HandleCancel completed successfully for user %d", userID)
	return nil
}

// isValidImageFormat validates image MIME type (T026 - FR-003)
// Accepts: image/jpeg, image/png, image/webp
func isValidImageFormat(mimeType string) bool {
	validFormats := []string{"image/jpeg", "image/png", "image/webp"}
	for _, valid := range validFormats {
		if mimeType == valid {
			return true
		}
	}
	return false
}

// sendError sends an error message to the user
// Helper for T031, T032, T033 error handling
func (h *EstimateHandler) sendError(c telebot.Context, message string) error {
	_, err := c.Bot().Send(c.Sender(), "❌ "+message)
	return err
}
