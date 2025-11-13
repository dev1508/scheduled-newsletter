package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HTTPConfig holds HTTP email API configuration
type HTTPConfig struct {
	APIKey    string
	FromEmail string
	FromName  string
	BaseURL   string
}

// BrevoEmailRequest represents Brevo API email request
type BrevoEmailRequest struct {
	Sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HTMLContent string `json:"htmlContent,omitempty"`
	TextContent string `json:"textContent,omitempty"`
}

// HTTPEmailSender handles HTTP-based email sending via Brevo API
type HTTPEmailSender struct {
	config *HTTPConfig
	client *http.Client
	logger *zap.Logger
}

// NewHTTPEmailSender creates a new HTTP email sender
func NewHTTPEmailSender(config *HTTPConfig, logger *zap.Logger) *HTTPEmailSender {
	return &HTTPEmailSender{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Send sends an email via Brevo HTTP API
func (h *HTTPEmailSender) Send(req *EmailRequest) error {
	// Prepare Brevo API request
	brevoReq := BrevoEmailRequest{
		Subject:     req.Subject,
		HTMLContent: req.HTMLBody,
		TextContent: req.TextBody,
	}

	// Set sender
	brevoReq.Sender.Name = h.config.FromName
	brevoReq.Sender.Email = h.config.FromEmail

	// Set recipient
	brevoReq.To = []struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}{
		{
			Email: req.To,
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(brevoReq)
	if err != nil {
		h.logger.Error("Failed to marshal email request", zap.Error(err))
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", h.config.BaseURL+"/v3/smtp/email", bytes.NewBuffer(jsonData))
	if err != nil {
		h.logger.Error("Failed to create HTTP request", zap.Error(err))
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("api-key", h.config.APIKey)

	// Send request
	resp, err := h.client.Do(httpReq)
	if err != nil {
		h.logger.Error("Failed to send HTTP request", zap.Error(err))
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		h.logger.Error("Brevo API returned error status",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status),
		)
		return fmt.Errorf("brevo API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	h.logger.Info("Email sent successfully via Brevo HTTP API",
		zap.String("to", req.To),
		zap.String("subject", req.Subject),
		zap.Int("status_code", resp.StatusCode),
	)

	return nil
}
