package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"api/biz/say_right/service"

	"github.com/gin-gonic/gin"
)

type PaddleHandler struct {
	svc service.UserService
}

func NewPaddleHandler(svc service.UserService) *PaddleHandler {
	return &PaddleHandler{
		svc: svc,
	}
}

// Webhook payload structure (simplified)
type PaddleEvent struct {
	EventID    string                 `json:"event_id"`
	EventType  string                 `json:"event_type"`
	OccurredAt string                 `json:"occurred_at"`
	Data       map[string]interface{} `json:"data"`
}

func (h *PaddleHandler) HandleWebhook(c *gin.Context) {
	// 1. Verify Signature
	signature := c.GetHeader("Paddle-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
		return
	}

	// Verify signature
	secret := os.Getenv("PADDLE_WEBHOOK_SECRET_KEY")
	if secret != "" {
		if !verifyPaddleSignature(signature, body, secret) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
			return
		}
	} else {
		// Warning: running without signature verification
		// log.Println("Warning: PADDLE_WEBHOOK_SECRET_KEY not set, skipping signature verification")
	}

	// 2. Parse Event
	var event PaddleEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 3. Handle Event
	switch event.EventType {
	case "transaction.completed":
		h.handleTransactionCompleted(c, event.Data)
	case "subscription.created":
		h.handleSubscriptionCreated(c, event.Data)
	case "subscription.updated":
		h.handleSubscriptionUpdated(c, event.Data)
	default:
		// Ignore other events
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *PaddleHandler) handleTransactionCompleted(c *gin.Context, data map[string]interface{}) {
	// Try to get email from custom_data
	if customData, ok := data["custom_data"].(map[string]interface{}); ok {
		if email, ok := customData["email"].(string); ok && email != "" {
			if err := h.svc.UpgradeUserToPro(c.Request.Context(), email); err != nil {
				// Log error (in a real app, use a logger)
				// fmt.Printf("Failed to upgrade user %s: %v\n", email, err)
			}
			return
		}
	}

	// TODO: Handle case where email is not in custom_data (e.g. check customer_id and lookup)
}

func (h *PaddleHandler) handleSubscriptionCreated(c *gin.Context, data map[string]interface{}) {
	// TODO: Update user subscription status
}

func (h *PaddleHandler) handleSubscriptionUpdated(c *gin.Context, data map[string]interface{}) {
	// TODO: Handle updates (renewals, cancellations)
}

// Verify signature helper (placeholder implementation)
// See https://developer.paddle.com/webhook-reference/verifying-webhooks
func verifyPaddleSignature(signatureHeader string, body []byte, secret string) bool {
	// Parse ts and h1 from header "ts=...;h1=..."
	parts := strings.Split(signatureHeader, ";")
	var ts, h1 string
	for _, part := range parts {
		if strings.HasPrefix(part, "ts=") {
			ts = strings.TrimPrefix(part, "ts=")
		} else if strings.HasPrefix(part, "h1=") {
			h1 = strings.TrimPrefix(part, "h1=")
		}
	}

	if ts == "" || h1 == "" {
		return false
	}

	// Prevent replay attacks (e.g., check if ts is within tolerance)
	// timestamp, _ := strconv.ParseInt(ts, 10, 64)
	// if time.Now().Unix() - timestamp > 60 { return false }

	payload := ts + ":" + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedMac := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(h1), []byte(expectedMac))
}
