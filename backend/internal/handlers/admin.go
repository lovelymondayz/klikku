package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/config"
	"klikku/internal/models"
	"klikku/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func AdminListMerchants(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(),
			"SELECT id, business_name, slug, logo_url, primary_color, subscription_status, created_at FROM merchants ORDER BY created_at DESC")
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch merchants")
		}
		defer rows.Close()

		var merchants []map[string]interface{}
		for rows.Next() {
			var id, businessName, slug, logoURL, primaryColor, subscriptionStatus string
			var createdAt time.Time
			err := rows.Scan(&id, &businessName, &slug, &logoURL, &primaryColor, &subscriptionStatus, &createdAt)
			if err != nil {
				continue
			}
			merchants = append(merchants, map[string]interface{}{
				"id":                 id,
				"business_name":      businessName,
				"slug":               slug,
				"logo_url":           logoURL,
				"primary_color":      primaryColor,
				"subscription_status": subscriptionStatus,
				"created_at":         createdAt,
			})
		}

		return utils.Success(c, merchants)
	}
}

func AdminCreateMerchant(db *pgxpool.Pool, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			BusinessName string `json:"business_name"`
			Slug         string `json:"slug"`
			AdminName    string `json:"admin_name"`
			AdminEmail   string `json:"admin_email"`
			AdminPassword string `json:"admin_password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		if req.BusinessName == "" || req.AdminEmail == "" || req.AdminPassword == "" {
			return utils.Error(c, fiber.StatusBadRequest, "missing required fields")
		}

		slug := req.Slug
		if slug == "" {
			slug = req.AdminEmail
		}

		// Create merchant
		var merchantID string
		err := db.QueryRow(context.Background(),
			"INSERT INTO merchants (business_name, slug, subscription_status) VALUES ($1, $2, 'active') RETURNING id",
			req.BusinessName, slug).Scan(&merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create merchant")
		}

		// Hash password and create admin user
		hash, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to process password")
		}

		var userID string
		name := req.AdminName
		if name == "" {
			name = "Admin"
		}
		err = db.QueryRow(context.Background(),
			"INSERT INTO users (name, email, password_hash, role, merchant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			name, req.AdminEmail, string(hash), models.RoleMerchantAdmin, merchantID).Scan(&userID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create admin user")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "merchant_id": merchantID, "admin_id": userID})
	}
}

func AdminUpdateMerchant(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid merchant ID")
		}

		var req struct {
			BusinessName     string `json:"business_name"`
			SubscriptionStatus string `json:"subscription_status"`
			PrimaryColor     string `json:"primary_color"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		_, err = db.Exec(context.Background(),
			"UPDATE merchants SET business_name = COALESCE(NULLIF($1, ''), business_name), subscription_status = COALESCE(NULLIF($2, ''), subscription_status), primary_color = COALESCE(NULLIF($3, ''), primary_color) WHERE id = $4",
			req.BusinessName, req.SubscriptionStatus, req.PrimaryColor, id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to update merchant")
		}

		return utils.Message(c, "merchant updated")
	}
}

func AdminDeleteMerchant(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid merchant ID")
		}

		_, err = db.Exec(context.Background(), "DELETE FROM merchants WHERE id = $1", id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to delete merchant")
		}

		return utils.Message(c, "merchant deleted")
	}
}

func AdminListAllSessions(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(),
			`SELECT s.id, s.merchant_id, s.status, s.payment_status, s.email, s.created_at, m.business_name 
			 FROM photobooth_sessions s
			 JOIN merchants m ON m.id = s.merchant_id
			 ORDER BY s.created_at DESC LIMIT 100`)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch sessions")
		}
		defer rows.Close()

		var sessions []map[string]interface{}
		for rows.Next() {
			var id, merchantID, status, paymentStatus, email, businessName string
			var createdAt time.Time
			err := rows.Scan(&id, &merchantID, &status, &paymentStatus, &email, &createdAt, &businessName)
			if err != nil {
				continue
			}
			sessions = append(sessions, map[string]interface{}{
				"id":             id,
				"merchant_id":    merchantID,
				"status":         status,
				"payment_status": paymentStatus,
				"email":          email,
				"created_at":     createdAt,
				"business_name":  businessName,
			})
		}

		return utils.Success(c, sessions)
	}
}

func AdminGetPlatformAnalytics(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var merchants, activeDevices, totalSessions, photos, prints, emails int
		var revenue float64

		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM merchants").Scan(&merchants)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM devices WHERE status = 'online'").Scan(&activeDevices)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM photobooth_sessions").Scan(&totalSessions)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM photos").Scan(&photos)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM print_jobs").Scan(&prints)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM email_deliveries").Scan(&emails)
		db.QueryRow(context.Background(), "SELECT COALESCE(SUM(amount), 0) FROM payments").Scan(&revenue)

		return utils.Success(c, map[string]interface{}{
			"merchants":     merchants,
			"active_devices": activeDevices,
			"total_sessions": totalSessions,
			"photos":        photos,
			"prints":        prints,
			"emails":        emails,
			"revenue":       revenue,
		})
	}
}
