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

// FinalizeSession composes photos and generates the final image
func FinalizeSession(db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, err := uuid.Parse(sessionID)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		// Get session details
		var merchantID, templateID string
		err = db.QueryRow(context.Background(),
			"SELECT merchant_id, template_id FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&merchantID, &templateID)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		// Get photos
		rows, err := db.Query(context.Background(),
			"SELECT original_url FROM photos WHERE session_id = $1 ORDER BY position",
			sessionID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch photos")
		}
		defer rows.Close()

		var photoURLs []string
		for rows.Next() {
			var url string
			rows.Scan(&url)
			photoURLs = append(photoURLs, url)
		}

		if len(photoURLs) == 0 {
			return utils.Error(c, fiber.StatusBadRequest, "no photos found for session")
		}

		// Download photos to temp files
		var photoPaths []string
		tempDir := filepath.Join(os.TempDir(), "klikku", sessionID)
		os.MkdirAll(tempDir, 0755)

		for i, url := range photoURLs {
			data, err := storage.Download("originals", url)
			if err != nil {
				log.Printf("Failed to download photo %s: %v", url, err)
				continue
			}

			ext := filepath.Ext(url)
			if ext == "" {
				ext = ".jpg"
			}
			tempPath := filepath.Join(tempDir, fmt.Sprintf("photo_%d%s", i, ext))
			if err := os.WriteFile(tempPath, data, 0644); err == nil {
				photoPaths = append(photoPaths, tempPath)
			}
		}

		if len(photoPaths) == 0 {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to download photos")
		}

		// Get template layout config
		var layoutConfig map[string]interface{}
		err = db.QueryRow(context.Background(),
			"SELECT layout_config FROM templates WHERE id = $1", templateID).Scan(&layoutConfig)
		if err != nil || layoutConfig == nil {
			layoutConfig = map[string]interface{}{
				"output_width":  1200,
				"output_height": 1800,
			}
		}

		// Compose photos
		finalName := sessionID + "_final.jpg"
		finalPath := filepath.Join(tempDir, "final.jpg")

		if err := utils.ComposePhotos(photoPaths, layoutConfig, finalPath); err != nil {
			log.Printf("Compose failed: %v", err)
			// If composition fails, just use the first photo as the final
			finalPath = photoPaths[0]
		}

		// Read composed image
		finalData, err := os.ReadFile(finalPath)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to read final image")
		}

		// Upload to finals bucket
		err = storage.Upload("finals", finalName, finalData, "image/jpeg")
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to save final image")
		}

		// Update session
		_, err = db.Exec(context.Background(),
			"UPDATE photobooth_sessions SET status = 'COMPLETED', final_image_url = $1, completed_at = NOW() WHERE id = $2",
			finalName, sessionID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to update session")
		}

		// Cleanup
		os.RemoveAll(tempDir)

		return utils.Success(c, fiber.Map{
			"session_id":      sessionID,
			"final_image_url": finalName,
			"download_url":    "/api/download/" + sessionID,
		})
	}
}

// DownloadFinal serves the final composed image
func DownloadFinal(db *pgxpool.Pool, storage *utils.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, err := uuid.Parse(sessionID)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		var finalImageURL string
		err = db.QueryRow(context.Background(),
			"SELECT final_image_url FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		if finalImageURL == "" {
			return utils.Error(c, fiber.StatusNotFound, "final image not ready")
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
