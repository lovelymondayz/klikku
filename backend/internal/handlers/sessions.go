package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/utils"
)

func getMerchantID(c *gin.Context) string {
	role, _ := c.Get("role")
	if role == "SUPER_ADMIN" {
		return c.Query("merchant_id")
	}
	merchantID, _ := c.Get("merchant_id")
	if merchantID == nil {
		return ""
	}
	return merchantID.(string)
}

func GetAttractScreen(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceToken := c.Param("token")

		var merchantID, businessName, logoURL, primaryColor, welcomeMessage string
		var deviceName, currentCampaignID string

		err := db.QueryRow(context.Background(), `
			SELECT m.id, m.business_name, m.logo_url, m.primary_color, m.welcome_message, d.name, d.current_campaign_id
			FROM devices d
			JOIN merchants m ON m.id = d.merchant_id
			WHERE d.device_token = $1
		`, deviceToken).Scan(&merchantID, &businessName, &logoURL, &primaryColor, &welcomeMessage, &deviceName, &currentCampaignID)
		if err != nil {
			utils.Error(c, 404, "device not found")
			return
		}

		utils.Success(c, map[string]interface{}{
			"merchant_id":     merchantID,
			"business_name":   businessName,
			"logo_url":        logoURL,
			"primary_color":   primaryColor,
			"welcome_message": welcomeMessage,
			"device_name":     deviceName,
			"campaign_id":     currentCampaignID,
		})
	}
}

func CreateSession(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceToken := c.Param("token")

		var merchantID, deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT id, merchant_id FROM devices WHERE device_token = $1",
			deviceToken).Scan(&deviceID, &merchantID)
		if err != nil {
			utils.Error(c, 404, "device not found")
			return
		}

		var campaignID string
		err = db.QueryRow(context.Background(),
			"SELECT current_campaign_id FROM devices WHERE id = $1",
			deviceID).Scan(&campaignID)
		if err != nil {
			utils.Error(c, 500, "failed to get device campaign")
			return
		}

		sessionID := uuid.New().String()
		_, err = db.Exec(context.Background(),
			"INSERT INTO photobooth_sessions (id, merchant_id, device_id, campaign_id, status, payment_status) VALUES ($1, $2, $3, $4, 'STARTED', 'PENDING')",
			sessionID, merchantID, deviceID, campaignID)
		if err != nil {
			utils.Error(c, 500, "failed to create session")
			return
		}

		c.JSON(201, gin.H{"success": true, "session_id": sessionID})
	}
}

func CapturePhoto(db *pgxpool.Pool, storage *utils.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		form, err := c.MultipartForm()
		if err != nil {
			utils.Error(c, 400, "invalid form data")
			return
		}

		files := form.File["photos"]
		if len(files) == 0 {
			utils.Error(c, 400, "no photos uploaded")
			return
		}

		var photoURLs []string
		for i, file := range files {
			f, err := file.Open()
			if err != nil {
				continue
			}

			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				continue
			}

			objectName := sessionID + "/original_" + string(rune('0'+i)) + ".jpg"
			err = storage.Upload("originals", objectName, data, "image/jpeg")
			if err != nil {
				continue
			}

			photoID := uuid.New().String()
			_, err = db.Exec(context.Background(),
				"INSERT INTO photos (id, session_id, original_url, position) VALUES ($1, $2, $3, $4)",
				photoID, sessionID, objectName, i)
			if err == nil {
				photoURLs = append(photoURLs, objectName)
			}
		}

		_, _ = db.Exec(context.Background(),
			"UPDATE photobooth_sessions SET status = 'CAPTURING' WHERE id = $1",
			sessionID)

		utils.Success(c, gin.H{
			"session_id": sessionID,
			"photos":     photoURLs,
			"count":      len(photoURLs),
		})
	}
}

func GetSession(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var finalImageURL, status string
		err := db.QueryRow(context.Background(),
			"SELECT final_image_url, status FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL, &status)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		utils.Success(c, gin.H{
			"session_id":      sessionID,
			"status":          status,
			"final_image_url": finalImageURL,
		})
	}
}

func DownloadSession(db *pgxpool.Pool) gin.HandlerFunc {
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

		utils.Success(c, gin.H{"download_url": finalImageURL})
	}
}

func ListSessions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		rows, err := db.Query(context.Background(),
			`SELECT id, merchant_id, device_id, campaign_id, template_id, status, payment_status, email, final_image_url, created_at, completed_at 
			 FROM photobooth_sessions WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT 50`,
			merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch sessions")
			return
		}
		defer rows.Close()

		var sessions []map[string]interface{}
		for rows.Next() {
			var id, mID, deviceID, campaignID, templateID, status, paymentStatus, email string
			var finalImageURL string
			var createdAt time.Time
			var completedAt *time.Time
			err := rows.Scan(&id, &mID, &deviceID, &campaignID, &templateID, &status, &paymentStatus, &email, &finalImageURL, &createdAt, &completedAt)
			if err != nil {
				continue
			}
			sessions = append(sessions, map[string]interface{}{
				"id":             id,
				"merchant_id":    mID,
				"device_id":      deviceID,
				"campaign_id":    campaignID,
				"template_id":    templateID,
				"status":         status,
				"payment_status": paymentStatus,
				"email":          email,
				"final_image_url": finalImageURL,
				"created_at":     createdAt,
				"completed_at":   completedAt,
			})
		}

		utils.Success(c, sessions)
	}
}

func GetSessionDetail(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var mID, deviceID, campaignID, templateID, status, paymentStatus, email string
		var finalImageURL string
		var createdAt time.Time
		var completedAt *time.Time

		err := db.QueryRow(context.Background(),
			`SELECT id, merchant_id, device_id, campaign_id, template_id, status, payment_status, email, final_image_url, created_at, completed_at 
			 FROM photobooth_sessions WHERE id = $1`, id).Scan(
			&id, &mID, &deviceID, &campaignID, &templateID, &status, &paymentStatus, &email, &finalImageURL, &createdAt, &completedAt)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		session := map[string]interface{}{
			"id":             id,
			"merchant_id":    mID,
			"device_id":      deviceID,
			"campaign_id":    campaignID,
			"template_id":    templateID,
			"status":         status,
			"payment_status": paymentStatus,
			"email":          email,
			"final_image_url": finalImageURL,
			"created_at":     createdAt,
			"completed_at":   completedAt,
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, original_url, processed_url, final_url, position FROM photos WHERE session_id = $1 ORDER BY position",
			id)
		if err != nil {
			utils.Error(c, 500, "failed to fetch photos")
			return
		}
		defer rows.Close()

		var photos []map[string]interface{}
		for rows.Next() {
			var photoID, originalURL, processedURL, finalURL string
			var position int
			rows.Scan(&photoID, &originalURL, &processedURL, &finalURL, &position)
			photos = append(photos, map[string]interface{}{
				"id":            photoID,
				"original_url":  originalURL,
				"processed_url": processedURL,
				"final_url":     finalURL,
				"position":      position,
			})
		}

		utils.Success(c, map[string]interface{}{
			"session": session,
			"photos":  photos,
		})
	}
}

func DeleteSession(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(context.Background(), "DELETE FROM photobooth_sessions WHERE id = $1", id)
		if err != nil {
			utils.Error(c, 500, "failed to delete session")
			return
		}

		utils.Message(c, "session deleted")
	}
}

func GetBranding(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		var businessName, logoURL, primaryColor, secondaryColor, font, welcomeMessage string
		var idleBackgroundURL string
		var emailDesign, socialLinks map[string]interface{}

		err := db.QueryRow(context.Background(),
			"SELECT business_name, logo_url, primary_color, secondary_color, font, welcome_message, idle_background_url, email_design, social_links FROM merchants WHERE id = $1",
			merchantID).Scan(&businessName, &logoURL, &primaryColor, &secondaryColor, &font, &welcomeMessage, &idleBackgroundURL, &emailDesign, &socialLinks)
		if err != nil {
			utils.Error(c, 404, "merchant not found")
			return
		}

		utils.Success(c, map[string]interface{}{
			"business_name":      businessName,
			"logo_url":           logoURL,
			"primary_color":      primaryColor,
			"secondary_color":    secondaryColor,
			"font":               font,
			"welcome_message":    welcomeMessage,
			"idle_background_url": idleBackgroundURL,
			"email_design":       emailDesign,
			"social_links":       socialLinks,
		})
	}
}

func UpdateBranding(db *pgxpool.Pool, storage *utils.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		var req struct {
			BusinessName      string                 `json:"business_name"`
			LogoURL           string                 `json:"logo_url"`
			PrimaryColor      string                 `json:"primary_color"`
			SecondaryColor    string                 `json:"secondary_color"`
			Font              string                 `json:"font"`
			WelcomeMessage    string                 `json:"welcome_message"`
			IdleBackgroundURL string                 `json:"idle_background_url"`
			EmailDesign       map[string]interface{} `json:"email_design"`
			SocialLinks       map[string]interface{} `json:"social_links"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
		}

		_, err := db.Exec(context.Background(),
			`UPDATE merchants SET 
				business_name = COALESCE(NULLIF($1, ''), business_name),
				logo_url = COALESCE(NULLIF($2, ''), logo_url),
				primary_color = COALESCE(NULLIF($3, ''), primary_color),
				secondary_color = COALESCE(NULLIF($4, ''), secondary_color),
				font = COALESCE(NULLIF($5, ''), font),
				welcome_message = COALESCE(NULLIF($6, ''), welcome_message),
				idle_background_url = COALESCE(NULLIF($7, ''), idle_background_url),
				email_design = COALESCE($8, email_design),
				social_links = COALESCE($9, social_links)
			 WHERE id = $10`,
			req.BusinessName, req.LogoURL, req.PrimaryColor, req.SecondaryColor, req.Font, req.WelcomeMessage, req.IdleBackgroundURL, req.EmailDesign, req.SocialLinks, merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to update branding")
			return
		}

		utils.Message(c, "branding updated")
	}
}

func ListDevices(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, merchant_id, name, device_token, status, current_campaign_id, printer_config, last_seen_at, created_at FROM devices WHERE merchant_id = $1 ORDER BY created_at DESC",
			merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch devices")
			return
		}
		defer rows.Close()

		var devices []map[string]interface{}
		for rows.Next() {
			var id, mID, name, deviceToken, status, currentCampaignID string
			var printerConfig map[string]interface{}
			var lastSeenAt, createdAt time.Time
			err := rows.Scan(&id, &mID, &name, &deviceToken, &status, &currentCampaignID, &printerConfig, &lastSeenAt, &createdAt)
			if err != nil {
				continue
			}
			devices = append(devices, map[string]interface{}{
				"id":                  id,
				"merchant_id":         mID,
				"name":                name,
				"device_token":        deviceToken,
				"status":              status,
				"current_campaign_id": currentCampaignID,
				"printer_config":      printerConfig,
				"last_seen_at":        lastSeenAt,
				"created_at":          createdAt,
			})
		}

		utils.Success(c, devices)
	}
}

func CreateDevice(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		var req struct {
			Name          string                 `json:"name" binding:"required"`
			DeviceToken   string                 `json:"device_token"`
			PrinterConfig map[string]interface{} `json:"printer_config"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
		}

		deviceToken := req.DeviceToken
		if deviceToken == "" {
			deviceToken = uuid.New().String()
		}

		printerConfig := req.PrinterConfig
		if printerConfig == nil {
			printerConfig = make(map[string]interface{})
		}

		var id string
		err := db.QueryRow(context.Background(),
			"INSERT INTO devices (merchant_id, name, device_token, printer_config) VALUES ($1, $2, $3, $4) RETURNING id",
			merchantID, req.Name, deviceToken, printerConfig).Scan(&id)
		if err != nil {
			utils.Error(c, 500, "failed to create device")
			return
		}

		c.JSON(201, gin.H{"success": true, "id": id, "device_token": deviceToken})
	}
}

func GetDevice(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var mID, name, deviceToken, status, currentCampaignID string
		var printerConfig map[string]interface{}
		var lastSeenAt, createdAt time.Time

		err := db.QueryRow(context.Background(),
			"SELECT id, merchant_id, name, device_token, status, current_campaign_id, printer_config, last_seen_at, created_at FROM devices WHERE id = $1",
			id).Scan(&id, &mID, &name, &deviceToken, &status, &currentCampaignID, &printerConfig, &lastSeenAt, &createdAt)
		if err != nil {
			utils.Error(c, 404, "device not found")
			return
		}

		utils.Success(c, map[string]interface{}{
			"id":                  id,
			"merchant_id":         mID,
			"name":                name,
			"device_token":        deviceToken,
			"status":              status,
			"current_campaign_id": currentCampaignID,
			"printer_config":      printerConfig,
			"last_seen_at":        lastSeenAt,
			"created_at":          createdAt,
		})
	}
}

func UpdateDevice(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Name              string                 `json:"name"`
			Status            string                 `json:"status"`
			CurrentCampaignID string                 `json:"current_campaign_id"`
			PrinterConfig     map[string]interface{} `json:"printer_config"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
		}

		_, err := db.Exec(context.Background(),
			`UPDATE devices SET 
				name = COALESCE(NULLIF($1, ''), name),
				status = COALESCE(NULLIF($2, ''), status),
				current_campaign_id = COALESCE(NULLIF($3, ''), current_campaign_id),
				printer_config = COALESCE($4, printer_config)
			 WHERE id = $5`,
			req.Name, req.Status, req.CurrentCampaignID, req.PrinterConfig, id)
		if err != nil {
			utils.Error(c, 500, "failed to update device")
			return
		}

		utils.Message(c, "device updated")
	}
}

func DeleteDevice(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(context.Background(), "DELETE FROM devices WHERE id = $1", id)
		if err != nil {
			utils.Error(c, 500, "failed to delete device")
			return
		}

		utils.Message(c, "device deleted")
	}
}

func GetAnalyticsOverview(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		var sessions, photos, prints, emails int
		var revenue float64

		db.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM photobooth_sessions WHERE merchant_id = $1",
			merchantID).Scan(&sessions)
		db.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM photos p JOIN photobooth_sessions s ON p.session_id = s.id WHERE s.merchant_id = $1",
			merchantID).Scan(&photos)
		db.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM print_jobs WHERE device_id IN (SELECT id FROM devices WHERE merchant_id = $1)",
			merchantID).Scan(&prints)
		db.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM email_deliveries WHERE session_id IN (SELECT id FROM photobooth_sessions WHERE merchant_id = $1)",
			merchantID).Scan(&emails)
		db.QueryRow(context.Background(),
			"SELECT COALESCE(SUM(amount), 0) FROM payments WHERE session_id IN (SELECT id FROM photobooth_sessions WHERE merchant_id = $1)",
			merchantID).Scan(&revenue)

		utils.Success(c, map[string]interface{}{
			"sessions": sessions,
			"photos":   photos,
			"prints":   prints,
			"emails":   emails,
			"revenue":  revenue,
		})
	}
}

func UploadFile(db *pgxpool.Pool, storage *utils.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			utils.Error(c, 400, "no file uploaded")
			return
		}

		f, err := file.Open()
		if err != nil {
			utils.Error(c, 500, "failed to open file")
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			utils.Error(c, 500, "failed to read file")
			return
		}

		bucket := c.PostForm("bucket")
		if bucket == "" {
			bucket = "assets"
		}
		objectName := uuid.New().String() + "_" + file.Filename

		err = storage.Upload(bucket, objectName, data, file.Header.Get("Content-Type"))
		if err != nil {
			utils.Error(c, 500, "failed to upload file")
			return
		}

		utils.Success(c, gin.H{
			"url":      objectName,
			"bucket":   bucket,
			"filename": file.Filename,
		})
	}
}

func FinalizeSession(db *pgxpool.Pool, storage *utils.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var merchantID, templateID string
		err := db.QueryRow(context.Background(),
			"SELECT merchant_id, template_id FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&merchantID, &templateID)
		if err != nil {
			utils.Error(c, 404, "session not found")
			return
		}

		rows, err := db.Query(context.Background(),
			"SELECT original_url FROM photos WHERE session_id = $1 ORDER BY position",
			sessionID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch photos")
			return
		}
		defer rows.Close()

		var photoURLs []string
		for rows.Next() {
			var url string
			rows.Scan(&url)
			photoURLs = append(photoURLs, url)
		}

		if len(photoURLs) == 0 {
			utils.Error(c, 400, "no photos found for session")
			return
		}

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
			utils.Error(c, 500, "failed to download photos")
			return
		}

		var layoutConfig map[string]interface{}
		err = db.QueryRow(context.Background(),
			"SELECT layout_config FROM templates WHERE id = $1", templateID).Scan(&layoutConfig)
		if err != nil || layoutConfig == nil {
			layoutConfig = map[string]interface{}{
				"output_width":  1200,
				"output_height": 1800,
			}
		}

		finalName := sessionID + "_final.jpg"
		finalPath := filepath.Join(tempDir, "final.jpg")

		if err := utils.ComposePhotos(photoPaths, layoutConfig, finalPath); err != nil {
			log.Printf("Compose failed: %v", err)
			finalPath = photoPaths[0]
		}

		finalData, err := os.ReadFile(finalPath)
		if err != nil {
			utils.Error(c, 500, "failed to read final image")
			return
		}

		err = storage.Upload("finals", finalName, finalData, "image/jpeg")
		if err != nil {
			utils.Error(c, 500, "failed to save final image")
			return
		}

		_, err = db.Exec(context.Background(),
			"UPDATE photobooth_sessions SET status = 'COMPLETED', final_image_url = $1, completed_at = NOW() WHERE id = $2",
			finalName, sessionID)
		if err != nil {
			utils.Error(c, 500, "failed to update session")
			return
		}

		os.RemoveAll(tempDir)

		utils.Success(c, gin.H{
			"session_id":      sessionID,
			"final_image_url": finalName,
			"download_url":    "/api/download/" + sessionID,
		})
	}
}

func DownloadFinal(db *pgxpool.Pool, storage *utils.Storage) gin.HandlerFunc {
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

		data, err := storage.Download("finals", finalImageURL)
		if err != nil {
			utils.Error(c, 404, "image not found")
			return
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Disposition", "inline; filename=\"photobooth.jpg\"")
		c.Data(200, "image/jpeg", data)
	}
}
