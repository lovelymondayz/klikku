package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/config"
	"klikku/internal/utils"
)

// SendEmailDelivery sends the photobooth final image to customer's email
func SendEmailDelivery(db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, err := uuid.Parse(sessionID)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		// Get session details with merchant info
		var finalImageURL, customerEmail, merchantName, merchantLogo string
		err = db.QueryRow(context.Background(), `
			SELECT s.final_image_url, s.email, m.business_name, m.logo_url
			FROM photobooth_sessions s
			JOIN merchants m ON m.id = s.merchant_id
			WHERE s.id = $1
		`, sessionID).Scan(&finalImageURL, &customerEmail, &merchantName, &merchantLogo)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		if finalImageURL == "" {
			return utils.Error(c, fiber.StatusBadRequest, "final image not ready")
		}

		if customerEmail == "" {
			// Try to get email from request body
			var req struct {
				Email string `json:"email"`
			}
			if err := c.BodyParser(&req); err == nil && req.Email != "" {
				customerEmail = req.Email
				// Update session with email
				db.Exec(context.Background(), "UPDATE photobooth_sessions SET email = $1 WHERE id = $2", customerEmail, sessionID)
			} else {
				return utils.Error(c, fiber.StatusBadRequest, "customer email required")
			}
		}

		// Download final image for email attachment
		finalData, err := storage.Download("finals", finalImageURL)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to load final image")
		}

		// Save to temp file for email
		tempDir := filepath.Join(os.TempDir(), "klikku", "email")
		os.MkdirAll(tempDir, 0755)
		tempPath := filepath.Join(tempDir, sessionID+"_final.jpg")
		if err := os.WriteFile(tempPath, finalData, 0644); err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to save temp image")
		}
		defer os.Remove(tempPath)

		// Generate email HTML
		downloadURL := fmt.Sprintf("%s/api/download/%s", cfg.FrontendURL, sessionID)
		htmlContent := utils.GenerateBrandedEmail(merchantName, merchantLogo, downloadURL, downloadURL)

		// Send email via Brevo
		brevo := utils.NewBrevoEmail(cfg.BrevoAPIKey, cfg.BrevoSender, cfg.BrevoSenderEmail)
		err = brevo.SendEmail(customerEmail, "", "Your photobooth memory is ready!", htmlContent)
		if err != nil {
			log.Printf("Email send failed: %v", err)
			// Record failed delivery
			db.Exec(context.Background(),
				"INSERT INTO email_deliveries (session_id, email, status) VALUES ($1, $2, 'FAILED')",
				sessionID, customerEmail)
			return utils.Error(c, fiber.StatusInternalServerError, "failed to send email")
		}

		// Record successful delivery
		_, err = db.Exec(context.Background(),
			"INSERT INTO email_deliveries (session_id, email, status, sent_at) VALUES ($1, $2, 'SENT', NOW())",
			sessionID, customerEmail)
		if err != nil {
			log.Printf("Failed to record email delivery: %v", err)
		}

		return utils.Message(c, "email sent successfully")
	}
}

// ResendEmail resends the email delivery for a session
func ResendEmail(db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Reuse same logic
		return SendEmailDelivery(db, storage, cfg)(c)
	}
}
