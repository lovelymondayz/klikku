package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/utils"
)

// CreatePrintJob creates a new print job for a session
func CreatePrintJob(db *pgxpool.Pool, storage *utils.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, parseErr := uuid.Parse(sessionID)
		if parseErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		var req struct {
			PrintType string `json:"print_type"`
			Copies    int    `json:"copies"`
		}
		_ = c.BodyParser(&req)

		printType := req.PrintType
		if printType == "" {
			printType = "4x6"
		}
		copies := req.Copies
		if copies == 0 {
			copies = 1
		}

		var deviceID string
		queryErr := db.QueryRow(context.Background(),
			"SELECT device_id FROM photobooth_sessions WHERE id = $1", sessionID).Scan(&deviceID)
		if queryErr != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		printJobID := uuid.New().String()
		_, execErr := db.Exec(context.Background(),
			"INSERT INTO print_jobs (id, session_id, device_id, print_type, copies, status) VALUES ($1, $2, $3, $4, $5, 'QUEUED')",
			printJobID, sessionID, deviceID, printType, copies)
		if execErr != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create print job")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":         printJobID,
			"session_id": sessionID,
			"status":     "QUEUED",
		})
	}
}

// GetPrintJob gets the status of a print job
func GetPrintJob(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid print job ID")
		}

		var sessionID, deviceID, printType, status, printerName, errorMessage string
		var copies int
		var createdAt, printedAt string

		err := db.QueryRow(context.Background(),
			"SELECT id, session_id, device_id, print_type, copies, status, created_at, printed_at FROM print_jobs WHERE id = $1",
			id).Scan(&id, &sessionID, &deviceID, &printType, &copies, &status, &printerName, &errorMessage, &createdAt, &printedAt)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "print job not found")
		}

		return utils.Success(c, fiber.Map{
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

// UpdatePrintJobStatus updates print job status
func UpdatePrintJobStatus(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid print job ID")
		}

		var req struct {
			Status       string `json:"status"`
			PrinterName  string `json:"printer_name"`
			ErrorMessage string `json:"error_message"`
		}
		if bodyErr := c.BodyParser(&req); bodyErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		if req.Status == "PRINTED" {
			_, err := db.Exec(context.Background(),
				"UPDATE print_jobs SET status = $1, printer_name = $2, printed_at = NOW() WHERE id = $3",
				req.Status, req.PrinterName, id)
			if err != nil {
				return utils.Error(c, fiber.StatusInternalServerError, "failed to update")
			}
		} else {
			_, err := db.Exec(context.Background(),
				"UPDATE print_jobs SET status = $1, printer_name = $2, error_message = $3 WHERE id = $4",
				req.Status, req.PrinterName, req.ErrorMessage, id)
			if err != nil {
				return utils.Error(c, fiber.StatusInternalServerError, "failed to update")
			}
		}

		return utils.Message(c, "print job updated")
	}
}

// ListPrintJobs lists all print jobs for a merchant
func ListPrintJobs(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		rows, err := db.Query(context.Background(),
			`SELECT p.id, p.session_id, p.device_id, p.print_type, p.copies, p.status, p.printer_name, p.error_message, p.created_at, p.printed_at
			 FROM print_jobs p
			 JOIN photobooth_sessions s ON s.id = p.session_id
			 WHERE s.merchant_id = $1
			 ORDER BY p.created_at DESC LIMIT 50`,
			merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch print jobs")
		}
		defer rows.Close()

		var jobs []fiber.Map
		for rows.Next() {
			var id, sessionID, deviceID, printType, status, printerName, errorMessage string
			var copies int
			var createdAt, printedAt string
			rows.Scan(&id, &sessionID, &deviceID, &printType, &copies, &status, &printerName, &errorMessage, &createdAt, &printedAt)
			jobs = append(jobs, fiber.Map{
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

		return utils.Success(c, jobs)
	}
}

// Reprint creates a reprint job from an existing session
func Reprint(db *pgxpool.Pool, storage *utils.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, parseErr := uuid.Parse(sessionID)
		if parseErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		var deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT device_id FROM photobooth_sessions WHERE id = $1", sessionID).Scan(&deviceID)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		printJobID := uuid.New().String()
		_, execErr := db.Exec(context.Background(),
			"INSERT INTO print_jobs (id, session_id, device_id, print_type, copies, status) VALUES ($1, $2, $3, '4x6', 1, 'QUEUED')",
			printJobID, sessionID, deviceID)
		if execErr != nil {
			log.Printf("Failed to create reprint job: %v", execErr)
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create reprint job")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":         printJobID,
			"session_id": sessionID,
			"status":     "QUEUED",
		})
	}
}

// GetPendingPrintJobs gets pending print jobs for a device
func GetPendingPrintJobs(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceID := c.Query("device_id")
		if deviceID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "device_id required")
		}

		rows, err := db.Query(context.Background(),
			`SELECT p.id, p.session_id, p.print_type, p.copies, s.final_image_url
			 FROM print_jobs p
			 JOIN photobooth_sessions s ON s.id = p.session_id
			 WHERE p.device_id = $1 AND p.status = 'QUEUED'
			 ORDER BY p.created_at LIMIT 10`,
			deviceID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch print jobs")
		}
		defer rows.Close()

		var jobs []fiber.Map
		for rows.Next() {
			var id, sessionID, printType, finalImageURL string
			var copies int
			rows.Scan(&id, &sessionID, &printType, &copies, &finalImageURL)
			jobs = append(jobs, fiber.Map{
				"id":               id,
				"session_id":       sessionID,
				"print_type":       printType,
				"copies":           copies,
				"final_image_url":  finalImageURL,
			})
		}

		return utils.Success(c, jobs)
	}
}

// AutoPrintJob automatically creates a print job after session completion
func AutoPrintJob(db *pgxpool.Pool, storage *utils.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, parseErr := uuid.Parse(sessionID)
		if parseErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		var deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT device_id FROM photobooth_sessions WHERE id = $1", sessionID).Scan(&deviceID)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		printJobID := uuid.New().String()
		_, execErr := db.Exec(context.Background(),
			"INSERT INTO print_jobs (id, session_id, device_id, print_type, copies, status) VALUES ($1, $2, $3, '4x6', 1, 'QUEUED')",
			printJobID, sessionID, deviceID)
		if execErr != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create print job")
		}

		// Start print flow simulation
		go simulatePrintFlow(db, printJobID)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":     printJobID,
			"status": "QUEUED",
		})
	}
}

// simulatePrintFlow simulates the print job lifecycle
func simulatePrintFlow(db *pgxpool.Pool, printJobID string) {
	states := []string{"PREPARING", "SENDING", "PRINTING", "PRINT_COMPLETE"}
	for _, state := range states {
		time.Sleep(2 * time.Second)
		_, err := db.Exec(context.Background(),
			"UPDATE print_jobs SET status = $1, printer_name = 'Default Printer' WHERE id = $2",
			state, printJobID)
		if err != nil {
			log.Printf("Print flow update failed: %v", err)
			return
		}
	}
}

// GenerateSecureDownloadURL creates a time-limited download token
func GenerateSecureDownloadURL(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, parseErr := uuid.Parse(sessionID)
		if parseErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		var finalImageURL string
		err := db.QueryRow(context.Background(),
			"SELECT final_image_url FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		if finalImageURL == "" {
			return utils.Error(c, fiber.StatusNotFound, "final image not ready")
		}

		// Generate expiring token (24 hours)
		token := uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)

		_, execErr := db.Exec(context.Background(),
			"INSERT INTO download_tokens (token, session_id, expires_at) VALUES ($1, $2, $3)",
			token, sessionID, expiresAt)
		if execErr != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to generate token")
		}

		return utils.Success(c, fiber.Map{
			"download_url": fmt.Sprintf("%s/api/download/%s/secure?token=%s", getBaseURL(), sessionID, token),
			"expires_at":   expiresAt,
		})
	}
}

// ValidateDownloadToken validates a download token
func ValidateDownloadToken(db *pgxpool.Pool, storage *utils.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		token := c.Query("token")

		if token == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "download token required")
		}

		var expiresAt time.Time
		err := db.QueryRow(context.Background(),
			"SELECT expires_at FROM download_tokens WHERE token = $1 AND session_id = $2",
			token, sessionID).Scan(&expiresAt)
		if err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "invalid token")
		}

		if time.Now().After(expiresAt) {
			return utils.Error(c, fiber.StatusUnauthorized, "download link expired")
		}

		// Serve image
		var finalImageURL string
		err = db.QueryRow(context.Background(),
			"SELECT final_image_url FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		data, err := storage.Download("finals", finalImageURL)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "image not found")
		}

		c.Set("Content-Type", "image/jpeg")
		c.Set("Content-Disposition", "inline; filename=\"photobooth.jpg\"")
		return c.Send(data)
	}
}

func getBaseURL() string {
	return "http://localhost:8083"
}
