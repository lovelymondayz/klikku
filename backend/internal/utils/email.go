package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BrevoEmail handles email sending via Brevo API
type BrevoEmail struct {
	APIKey    string
	Sender    string
	SenderEmail string
}

// NewBrevoEmail creates a new Brevo email client
func NewBrevoEmail(apiKey, sender, senderEmail string) *BrevoEmail {
	return &BrevoEmail{
		APIKey:      apiKey,
		Sender:      sender,
		SenderEmail: senderEmail,
	}
}

// EmailRequest represents the Brevo API request
type EmailRequest struct {
	Sender    EmailAddress `json:"sender"`
	To        []EmailAddress `json:"to"`
	Subject   string      `json:"subject"`
	HTMLContent string   `json:"htmlContent"`
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// SendEmail sends an email via Brevo API
func (b *BrevoEmail) SendEmail(toEmail, toName, subject, htmlContent string) error {
	if b.APIKey == "" {
		return fmt.Errorf("Brevo API key not configured")
	}

	reqBody := EmailRequest{
		Sender: EmailAddress{
			Email: b.SenderEmail,
			Name:  b.Sender,
		},
		To: []EmailAddress{
			{Email: toEmail, Name: toName},
		},
		Subject:     subject,
		HTMLContent: htmlContent,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal email request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", b.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send email request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("brevo API error (status %d): %s", resp.StatusCode, string(body))
}

// GenerateBrandedEmail creates the HTML email content for photobooth delivery
func GenerateBrandedEmail(merchantName, logoURL, finalImageURL, downloadURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0; padding:0; background-color:#f8f9fa; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
	<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="max-width:600px; margin:0 auto; background:#ffffff;">
		<tr>
			<td style="padding:40px 30px; text-align:center;">
				%s
				<h1 style="color:#333; font-size:24px; margin:0 0 10px;">%s</h1>
				<p style="color:#666; font-size:16px; margin:0 0 30px;">Your photobooth memory is ready!</p>
			</td>
		</tr>
		<tr>
			<td style="padding:0 30px; text-align:center;">
				<div style="background:linear-gradient(135deg,#ff6b9d,#c44dff); border-radius:20px; padding:4px; display:inline-block;">
					<img src="%s" alt="Your Photo" style="max-width:100%%; border-radius:16px; display:block;">
				</div>
			</td>
		</tr>
		<tr>
			<td style="padding:30px; text-align:center;">
				<a href="%s" style="display:inline-block; background:linear-gradient(135deg,#ff6b9d,#c44dff); color:#ffffff; text-decoration:none; padding:16px 40px; border-radius:50px; font-size:18px; font-weight:bold; box-shadow:0 4px 15px rgba(255,107,157,0.4);">
					Download Your Photo
				</a>
			</td>
		</tr>
		<tr>
			<td style="padding:20px 30px; text-align:center;">
				<p style="color:#999; font-size:12px;">Follow us for more moments ✨</p>
			</td>
		</tr>
	</table>
</body>
</html>`, logoImgTag(logoURL), merchantName, finalImageURL, downloadURL)
}

func logoImgTag(logoURL string) string {
	if logoURL == "" {
		return ""
	}
	return fmt.Sprintf(`<img src="%s" alt="Logo" style="width:80px; height:80px; object-fit:contain; border-radius:16px; margin-bottom:20px;">`, logoURL)
}
