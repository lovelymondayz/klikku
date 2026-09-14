package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/storage"
	"klikku/internal/utils"
)

func SendEmailDelivery(db *pgxpool.Pool, store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var finalImageURL, customerEmail, merchantName, merchantLogo string
		err := db.QueryRow(context.Background(), `
			SELECT s.final_image_url, s.email, m.business_name, m.logo_url
			FROM photobooth_sessions s
			JOIN merchants m ON m.id = s.merchant_id
			WHERE s.id = $1
		`, sessionID).Scan(&finalImageURL, &customerEmail, &merchantName, &merchantLogo)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		if finalImageURL == "" {
			utils.Error(c, 400, "final image not ready")
			return
		}

		if customerEmail == "" {
			var req struct {
				Email string `json:"email"`
			}
			if err := c.ShouldBindJSON(&req); err == nil && req.Email != "" {
				customerEmail = req.Email
				db.Exec(context.Background(), "UPDATE photobooth_sessions SET email = $1 WHERE id = $2", customerEmail, sessionID)
			} else {
				utils.Error(c, 400, "customer email required")
				return
			}
		}

		finalData, err := store.Download("finals", finalImageURL)
		if err != nil {
			utils.Error(c, 500, "failed to load final image")
			return
		}

		tempDir := filepath.Join(os.TempDir(), "klikku", "email")
		os.MkdirAll(tempDir, 0755)
		tempPath := filepath.Join(tempDir, sessionID+"_final.jpg")
		if err := os.WriteFile(tempPath, finalData, 0644); err != nil {
			utils.Error(c, 500, "failed to save temp image")
			return
		}
		defer os.Remove(tempPath)

		downloadURL := fmt.Sprintf("%s/api/download/%s", "https://klikku.arjism.com", sessionID)
		htmlContent := utils.GenerateBrandedEmail(merchantName, merchantLogo, downloadURL, downloadURL)

		brevo := utils.NewBrevoEmail("", "", "")
		err = brevo.SendEmail(customerEmail, "", "Your photobooth memory is ready!", htmlContent)
		if err != nil {
			log.Printf("Email send failed: %v", err)
			db.Exec(context.Background(),
				"INSERT INTO email_deliveries (session_id, email, status) VALUES ($1, $2, 'FAILED')",
				sessionID, customerEmail)
			utils.Error(c, 500, "failed to send email")
			return
		}

		_, err = db.Exec(context.Background(),
			"INSERT INTO email_deliveries (session_id, email, status, sent_at) VALUES ($1, $2, 'SENT', NOW())",
			sessionID, customerEmail)
		if err != nil {
			log.Printf("Failed to record email delivery: %v", err)
		}

		utils.Message(c, "email sent successfully")
	}
}

func ResendEmail(db *pgxpool.Pool, store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		SendEmailDelivery(db, store)(c)
	}
}

func CreatePrintJob(db *pgxpool.Pool, store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var req struct {
			PrintType string `json:"print_type"`
			Copies    int    `json:"copies"`
		}
		_ = c.ShouldBindJSON(&req)

		printType := req.PrintType
		if printType == "" {
			printType = "4x6"
		}
		copies := req.Copies
		if copies == 0 {
			copies = 1
		}

		var deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT device_id FROM photobooth_sessions WHERE id = $1", sessionID).Scan(&deviceID)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		printJobID := uuid.New().String()
		_, err = db.Exec(context.Background(),
			"INSERT INTO print_jobs (id, session_id, device_id, print_type, copies, status) VALUES ($1, $2, $3, $4, $5, 'QUEUED')",
			printJobID, sessionID, deviceID, printType, copies)
		if err != nil {
			utils.Error(c, 500, "failed to create print job")
			return
		}

		c.JSON(201, gin.H{
			"id":         printJobID,
			"session_id": sessionID,
			"status":     "QUEUED",
		})
	}
}

func GetPrintJob(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var sessionID, deviceID, printType, status, printerName, errorMessage string
		var copies int
		var createdAt, printedAt *time.Time

		err := db.QueryRow(context.Background(),
			"SELECT id, session_id, device_id, print_type, copies, status, printer_name, error_message, created_at, printed_at FROM print_jobs WHERE id = $1",
			id).Scan(&id, &sessionID, &deviceID, &printType, &copies, &status, &printerName, &errorMessage, &createdAt, &printedAt)
		if err != nil {
			utils.Error(c, 404, "print job not found")
			return
		}

		utils.Success(c, gin.H{
			"id":            id,
			"session_id":    sessionID,
			"device_id":     deviceID,
			"print_type":    printType,
			"copies":        copies,
			"status":        status,
			"printer_name":  printerName,
			"error_message": errorMessage,
			"created_at":    createdAt,
			"printed_at":    printedAt,
		})
	}
}

func UpdatePrintJobStatus(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Status       string `json:"status" binding:"required"`
			PrinterName  string `json:"printer_name"`
			ErrorMessage string `json:"error_message"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
		}

		if req.Status == "PRINTED" {
			_, err := db.Exec(context.Background(),
				"UPDATE print_jobs SET status = $1, printer_name = $2, printed_at = NOW() WHERE id = $3",
				req.Status, req.PrinterName, id)
			if err != nil {
				utils.Error(c, 500, "failed to update")
				return
			}
		} else {
			_, err := db.Exec(context.Background(),
				"UPDATE print_jobs SET status = $1, printer_name = $2, error_message = $3 WHERE id = $4",
				req.Status, req.PrinterName, req.ErrorMessage, id)
			if err != nil {
				utils.Error(c, 500, "failed to update")
				return
			}
		}

		utils.Message(c, "print job updated")
	}
}

func ListPrintJobs(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		rows, err := db.Query(context.Background(),
			`SELECT p.id, p.session_id, p.device_id, p.print_type, p.copies, p.status, p.printer_name, p.error_message, p.created_at, p.printed_at
			 FROM print_jobs p
			 JOIN photobooth_sessions s ON s.id = p.session_id
			 WHERE s.merchant_id = $1
			 ORDER BY p.created_at DESC LIMIT 50`,
			merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch print jobs")
			return
		}
		defer rows.Close()

		var jobs []gin.H
		for rows.Next() {
			var id, sessionID, deviceID, printType, status, printerName, errorMessage string
			var copies int
			var createdAt, printedAt *time.Time
			rows.Scan(&id, &sessionID, &deviceID, &printType, &copies, &status, &printerName, &errorMessage, &createdAt, &printedAt)
			jobs = append(jobs, gin.H{
				"id":            id,
				"session_id":    sessionID,
				"device_id":     deviceID,
				"print_type":    printType,
				"copies":        copies,
				"status":        status,
				"printer_name":  printerName,
				"error_message": errorMessage,
				"created_at":    createdAt,
				"printed_at":    printedAt,
			})
		}

		utils.Success(c, jobs)
	}
}

func Reprint(db *pgxpool.Pool, store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT device_id FROM photobooth_sessions WHERE id = $1", sessionID).Scan(&deviceID)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		printJobID := uuid.New().String()
		_, err = db.Exec(context.Background(),
			"INSERT INTO print_jobs (id, session_id, device_id, print_type, copies, status) VALUES ($1, $2, $3, '4x6', 1, 'QUEUED')",
			printJobID, sessionID, deviceID)
		if err != nil {
			log.Printf("Failed to create reprint job: %v", err)
			utils.Error(c, 500, "failed to create reprint job")
			return
		}

		c.JSON(201, gin.H{
			"id":         printJobID,
			"session_id": sessionID,
			"status":     "QUEUED",
		})
	}
}

func GetPendingPrintJobs(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Query("device_id")
		if deviceID == "" {
			utils.Error(c, 400, "device_id required")
			return
		}

		rows, err := db.Query(context.Background(),
			`SELECT p.id, p.session_id, p.print_type, p.copies, s.final_image_url
			 FROM print_jobs p
			 JOIN photobooth_sessions s ON s.id = p.session_id
			 WHERE p.device_id = $1 AND p.status = 'QUEUED'
			 ORDER BY p.created_at LIMIT 10`,
			deviceID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch print jobs")
			return
		}
		defer rows.Close()

		var jobs []gin.H
		for rows.Next() {
			var id, sessionID, printType, finalImageURL string
			var copies int
			rows.Scan(&id, &sessionID, &printType, &copies, &finalImageURL)
			jobs = append(jobs, gin.H{
				"id":              id,
				"session_id":      sessionID,
				"print_type":      printType,
				"copies":          copies,
				"final_image_url": finalImageURL,
			})
		}

		utils.Success(c, jobs)
	}
}

func AutoPrintJob(db *pgxpool.Pool, store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT device_id FROM photobooth_sessions WHERE id = $1", sessionID).Scan(&deviceID)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		printJobID := uuid.New().String()
		_, err = db.Exec(context.Background(),
			"INSERT INTO print_jobs (id, session_id, device_id, print_type, copies, status) VALUES ($1, $2, $3, '4x6', 1, 'QUEUED')",
			printJobID, sessionID, deviceID)
		if err != nil {
			utils.Error(c, 500, "failed to create print job")
			return
		}

		c.JSON(201, gin.H{
			"id":     printJobID,
			"status": "QUEUED",
		})
	}
}

func GenerateSecureDownloadURL(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var finalImageURL string
		err := db.QueryRow(context.Background(),
			"SELECT final_image_url FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		if finalImageURL == "" {
			utils.Error(c, 404, "final image not ready")
			return
		}

		token := uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)

		_, err = db.Exec(context.Background(),
			"INSERT INTO download_tokens (token, session_id, expires_at) VALUES ($1, $2, $3)",
			token, sessionID, expiresAt)
		if err != nil {
			utils.Error(c, 500, "failed to generate token")
			return
		}

		baseURL := "https://klikku.arjism.com"
		utils.Success(c, gin.H{
			"download_url": fmt.Sprintf("%s/api/download/%s/secure?token=%s", baseURL, sessionID, token),
			"expires_at":   expiresAt,
		})
	}
}

func ValidateDownloadToken(db *pgxpool.Pool, store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		token := c.Query("token")

		if token == "" {
			utils.Error(c, 401, "download token required")
			return
		}

		var expiresAt time.Time
		err := db.QueryRow(context.Background(),
			"SELECT expires_at FROM download_tokens WHERE token = $1 AND session_id = $2",
			token, sessionID).Scan(&expiresAt)
		if err != nil {
			utils.Error(c, 401, "invalid token")
			return
		}

		if time.Now().After(expiresAt) {
			utils.Error(c, 401, "download link expired")
			return
		}

		var finalImageURL string
		err = db.QueryRow(context.Background(),
			"SELECT final_image_url FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		data, err := store.Download("finals", finalImageURL)
		if err != nil {
			utils.Error(c, 404, "image not found")
			return
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Disposition", "inline; filename=\"photobooth.jpg\"")
		c.Data(200, "image/jpeg", data)
	}
}
