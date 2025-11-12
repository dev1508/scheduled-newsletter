package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

// EmailRequest represents an email to be sent
type EmailRequest struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
}

// SMTPSender handles SMTP email sending
type SMTPSender struct {
	config *SMTPConfig
	logger *zap.Logger
}

// NewSMTPSender creates a new SMTP sender
func NewSMTPSender(config *SMTPConfig, logger *zap.Logger) *SMTPSender {
	return &SMTPSender{
		config: config,
		logger: logger,
	}
}

// Send sends an email via SMTP
func (s *SMTPSender) Send(req *EmailRequest) error {
	// Create SMTP auth
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	
	// Build email message
	message := s.buildMessage(req)
	
	// SMTP server address
	addr := s.config.Host + ":" + s.config.Port
	
	// Send email using Go's built-in SMTP with STARTTLS
	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{req.To}, message)
	if err != nil {
		s.logger.Error("Failed to send email",
			zap.String("to", req.To),
			zap.String("subject", req.Subject),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email to %s: %w", req.To, err)
	}
	
	s.logger.Debug("Email sent successfully",
		zap.String("to", req.To),
		zap.String("subject", req.Subject),
	)
	
	return nil
}


// buildMessage constructs the email message
func (s *SMTPSender) buildMessage(req *EmailRequest) []byte {
	var message strings.Builder
	
	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail))
	message.WriteString(fmt.Sprintf("To: %s\r\n", req.To))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	
	// If both HTML and text body are provided, create multipart
	if req.HTMLBody != "" && req.TextBody != "" {
		boundary := "boundary-newsletter-email"
		message.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", boundary))
		message.WriteString("\r\n")
		
		// Text part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(req.TextBody)
		message.WriteString("\r\n")
		
		// HTML part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(req.HTMLBody)
		message.WriteString("\r\n")
		
		message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else if req.HTMLBody != "" {
		// HTML only
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(req.HTMLBody)
	} else {
		// Text only
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(req.TextBody)
	}
	
	return []byte(message.String())
}
