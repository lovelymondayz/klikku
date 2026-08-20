package handlers

import (
	"context"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/config"
	"klikku/internal/utils"
)

// Photobooth session handlers
func GetAttractScreen(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceToken := c.Params("token")

		var merchantID, businessName, logoURL, primaryColor, welcomeMessage string
		var deviceName, currentCampaignID string

		err := db.QueryRow(context.Background(), `
			SELECT m.id, m.business_name, m.logo_url, m.primary_color, m.welcome_message, d.name, d.current_campaign_id
			FROM devices d
			JOIN merchants m ON m.id = d.merchant_id
			WHERE d.device_token = $1
		`, deviceToken).Scan(&merchantID, &businessName, &logoURL, &primaryColor, &welcomeMessage, &deviceName, &currentCampaignID)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "device not found")
		}

		return utils.Success(c, map[string]interface{}{
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

func CreateSession(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceToken := c.Params("token")

		var merchantID, deviceID string
		err := db.QueryRow(context.Background(),
			"SELECT id, merchant_id FROM devices WHERE device_token = $1",
			deviceToken).Scan(&deviceID, &merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "device not found")
		}

		var campaignID string
		err = db.QueryRow(context.Background(),
			"SELECT current_campaign_id FROM devices WHERE id = $1",
			deviceID).Scan(&campaignID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to get device campaign")
		}

		sessionID := uuid.New().String()
		_, err = db.Exec(context.Background(),
			"INSERT INTO photobooth_sessions (id, merchant_id, device_id, campaign_id, status, payment_status) VALUES ($1, $2, $3, $4, 'STARTED', 'PENDING')",
			sessionID, merchantID, deviceID, campaignID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create session")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "session_id": sessionID})
	}
}

func CapturePhoto(db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		_, err := uuid.Parse(sessionID)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		form, err := c.MultipartForm()
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid form data")
		}

		files := form.File["photos"]
		if len(files) == 0 {
			return utils.Error(c, fiber.StatusBadRequest, "no photos uploaded")
		}

		var merchantID, templateID string
		err = db.QueryRow(context.Background(),
			"SELECT merchant_id, template_id FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&merchantID, &templateID)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		var photoURLs []string
		for i, file := range files {
			f, err := file.Open()
			if err != nil {
				continue
			}
			defer f.Close()

			data, err := io.ReadAll(f)
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

		return utils.Success(c, fiber.Map{
			"session_id": sessionID,
			"photos":     photoURLs,
			"count":      len(photoURLs),
		})
	}
}

func GetSession(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.Success(c, fiber.Map{"status": "ok"})
	}
}

func DownloadSession(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Params("id")

		var finalImageURL string
		err := db.QueryRow(context.Background(),
			"SELECT final_image_url FROM photobooth_sessions WHERE id = $1",
			sessionID).Scan(&finalImageURL)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		return utils.Success(c, fiber.Map{"download_url": finalImageURL})
	}
}

func ListSessions(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		rows, err := db.Query(context.Background(),
			`SELECT id, merchant_id, device_id, campaign_id, template_id, status, payment_status, email, final_image_url, created_at, completed_at 
			 FROM photobooth_sessions WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT 50`,
			merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch sessions")
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

		return utils.Success(c, sessions)
	}
}

func GetSessionDetail(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		var session map[string]interface{}
		var mID, deviceID, campaignID, templateID, status, paymentStatus, email string
		var finalImageURL string
		var createdAt time.Time
		var completedAt *time.Time

		err = db.QueryRow(context.Background(),
			`SELECT id, merchant_id, device_id, campaign_id, template_id, status, payment_status, email, final_image_url, created_at, completed_at 
			 FROM photobooth_sessions WHERE id = $1`, id).Scan(
			&id, &mID, &deviceID, &campaignID, &templateID, &status, &paymentStatus, &email, &finalImageURL, &createdAt, &completedAt)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "session not found")
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, original_url, processed_url, final_url, position FROM photos WHERE session_id = $1 ORDER BY position",
			id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch photos")
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

		return utils.Success(c, map[string]interface{}{
			"session": session,
			"photos":  photos,
		})
	}
}

func DeleteSession(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid session ID")
		}

		_, err = db.Exec(context.Background(), "DELETE FROM photobooth_sessions WHERE id = $1", id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to delete session")
		}

		return utils.Message(c, "session deleted")
	}
}

// Branding
func GetBranding(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		var businessName, logoURL, primaryColor, secondaryColor, font, welcomeMessage string
		var idleBackgroundURL string
		var emailDesign, socialLinks map[string]interface{}

		err := db.QueryRow(context.Background(),
			"SELECT business_name, logo_url, primary_color, secondary_color, font, welcome_message, idle_background_url, email_design, social_links FROM merchants WHERE id = $1",
			merchantID).Scan(&businessName, &logoURL, &primaryColor, &secondaryColor, &font, &welcomeMessage, &idleBackgroundURL, &emailDesign, &socialLinks)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "merchant not found")
		}

		return utils.Success(c, map[string]interface{}{
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

func UpdateBranding(db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		var req struct {
			BusinessName     string                 `json:"business_name"`
			LogoURL          string                 `json:"logo_url"`
			PrimaryColor     string                 `json:"primary_color"`
			SecondaryColor   string                 `json:"secondary_color"`
			Font             string                 `json:"font"`
			WelcomeMessage   string                 `json:"welcome_message"`
			IdleBackgroundURL string                `json:"idle_background_url"`
			EmailDesign      map[string]interface{} `json:"email_design"`
			SocialLinks      map[string]interface{} `json:"social_links"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
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
			return utils.Error(c, fiber.StatusInternalServerError, "failed to update branding")
		}

		return utils.Message(c, "branding updated")
	}
}

// Devices
func ListDevices(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, merchant_id, name, device_token, status, current_campaign_id, printer_config, last_seen_at, created_at FROM devices WHERE merchant_id = $1 ORDER BY created_at DESC",
			merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch devices")
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

		return utils.Success(c, devices)
	}
}

func CreateDevice(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		var req struct {
			Name        string                 `json:"name"`
			DeviceToken string                 `json:"device_token"`
			PrinterConfig map[string]interface{} `json:"printer_config"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		if req.Name == "" {
			return utils.Error(c, fiber.StatusBadRequest, "device name required")
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
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create device")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "id": id, "device_token": deviceToken})
	}
}

func GetDevice(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid device ID")
		}

		var mID, name, deviceToken, status, currentCampaignID string
		var printerConfig map[string]interface{}
		var lastSeenAt, createdAt time.Time

		err = db.QueryRow(context.Background(),
			"SELECT id, merchant_id, name, device_token, status, current_campaign_id, printer_config, last_seen_at, created_at FROM devices WHERE id = $1",
			id).Scan(&id, &mID, &name, &deviceToken, &status, &currentCampaignID, &printerConfig, &lastSeenAt, &createdAt)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "device not found")
		}

		return utils.Success(c, map[string]interface{}{
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

func UpdateDevice(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid device ID")
		}

		var req struct {
			Name              string                 `json:"name"`
			Status            string                 `json:"status"`
			CurrentCampaignID string                 `json:"current_campaign_id"`
			PrinterConfig     map[string]interface{} `json:"printer_config"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		_, err = db.Exec(context.Background(),
			`UPDATE devices SET 
				name = COALESCE(NULLIF($1, ''), name),
				status = COALESCE(NULLIF($2, ''), status),
				current_campaign_id = COALESCE(NULLIF($3, ''), current_campaign_id),
				printer_config = COALESCE($4, printer_config)
			 WHERE id = $5`,
			req.Name, req.Status, req.CurrentCampaignID, req.PrinterConfig, id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to update device")
		}

		return utils.Message(c, "device updated")
	}
}

func DeleteDevice(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid device ID")
		}

		_, err = db.Exec(context.Background(), "DELETE FROM devices WHERE id = $1", id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to delete device")
		}

		return utils.Message(c, "device deleted")
	}
}

// Analytics
func GetAnalyticsOverview(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
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

		return utils.Success(c, map[string]interface{}{
			"sessions": sessions,
			"photos":   photos,
			"prints":   prints,
			"emails":   emails,
			"revenue":  revenue,
		})
	}
}

// Upload file to local storage
func UploadFile(db *pgxpool.Pool, storage *utils.Storage, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "no file uploaded")
		}

		f, err := file.Open()
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to open file")
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to read file")
		}

		bucket := c.FormValue("bucket", "assets")
		objectName := uuid.New().String() + "_" + file.Filename

		err = storage.Upload(bucket, objectName, data, file.Header.Get("Content-Type"))
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to upload file")
		}

		return utils.Success(c, fiber.Map{
			"url":      objectName,
			"bucket":   bucket,
			"filename": file.Filename,
		})
	}
}
