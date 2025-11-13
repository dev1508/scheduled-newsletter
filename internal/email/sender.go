package email

import (
	"fmt"
	
	"go.uber.org/zap"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(req *EmailRequest) error
}

// UnifiedEmailSender can use either SMTP or HTTP API
type UnifiedEmailSender struct {
	smtpSender *SMTPSender
	httpSender *HTTPEmailSender
	useHTTP    bool
	logger     *zap.Logger
}

// UnifiedConfig holds configuration for both SMTP and HTTP
type UnifiedConfig struct {
	// SMTP Configuration
	SMTP *SMTPConfig
	
	// HTTP Configuration
	HTTP *HTTPConfig
	
	// Preference: true for HTTP, false for SMTP
	UseHTTP bool
}

// NewUnifiedEmailSender creates a new unified email sender
func NewUnifiedEmailSender(config *UnifiedConfig, logger *zap.Logger) *UnifiedEmailSender {
	sender := &UnifiedEmailSender{
		useHTTP: config.UseHTTP,
		logger:  logger,
	}

	// Initialize SMTP sender if config provided
	if config.SMTP != nil {
		sender.smtpSender = NewSMTPSender(config.SMTP, logger)
	}

	// Initialize HTTP sender if config provided
	if config.HTTP != nil {
		sender.httpSender = NewHTTPEmailSender(config.HTTP, logger)
	}

	return sender
}

// Send sends an email using the configured method (HTTP or SMTP)
func (u *UnifiedEmailSender) Send(req *EmailRequest) error {
	if u.useHTTP && u.httpSender != nil {
		u.logger.Debug("Sending email via HTTP API")
		return u.httpSender.Send(req)
	}

	if u.smtpSender != nil {
		u.logger.Debug("Sending email via SMTP")
		return u.smtpSender.Send(req)
	}

	u.logger.Error("No email sender configured")
	return fmt.Errorf("no email sender configured")
}
