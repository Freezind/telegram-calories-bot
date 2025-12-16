// Package handlers implements Telegram bot message and command handlers
package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/freezind/telegram-calories-bot/src/models"
	"github.com/freezind/telegram-calories-bot/src/services"
	telebot "gopkg.in/telebot.v3"
)

// EstimateHandler handles the /estimate command (T024)
// Initiates the calorie estimation flow per FR-001
type EstimateHandler struct {
	sessionManager *services.SessionManager
	geminiClient   *services.GeminiClient
}

// NewEstimateHandler creates a new EstimateHandler instance
func NewEstimateHandler(sm *services.SessionManager, gc *services.GeminiClient) *EstimateHandler {
	return &EstimateHandler{
		sessionManager: sm,
		geminiClient:   gc,
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
	btnReEstimate := markup.Data("🔄 Re-estimate", "re_estimate")
	btnCancel := markup.Data("❌ Cancel", "cancel")
	markup.Inline(
		markup.Row(btnReEstimate, btnCancel),
	)

	_, err = c.Bot().Send(c.Sender(), formattedResult, markup)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	// Keep session in AwaitingImage state for potential Re-estimate
	h.sessionManager.UpdateSession(userID, models.StateAwaitingImage)

	return nil
}

// HandleReEstimate handles the Re-estimate button click (User Story 2)
func (h *EstimateHandler) HandleReEstimate(c telebot.Context) error {
	userID := c.Sender().ID
	log.Printf("[HANDLER] HandleReEstimate called for user %d", userID)

	// Update callback to show feedback
	if err := c.Respond(&telebot.CallbackResponse{Text: "Send another image"}); err != nil {
		log.Printf("[HANDLER ERROR] Failed to respond to re_estimate callback for user %d: %v", userID, err)
	}

	// Delete previous result message
	if err := c.Delete(); err != nil {
		log.Printf("Failed to delete message: %v", err)
	}

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
func (h *EstimateHandler) HandleCancel(c telebot.Context) error {
	userID := c.Sender().ID
	log.Printf("[HANDLER] HandleCancel called for user %d", userID)

	// Update callback to show feedback
	if err := c.Respond(&telebot.CallbackResponse{Text: "Estimation canceled"}); err != nil {
		log.Printf("[HANDLER ERROR] Failed to respond to cancel callback for user %d: %v", userID, err)
	}

	// Delete the message with Cancel button
	if err := c.Delete(); err != nil {
		log.Printf("Failed to delete message: %v", err)
	}

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
