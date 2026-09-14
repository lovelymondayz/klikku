package handlers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/utils"
)

func AdminListMerchants(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(),
			"SELECT id, business_name, slug, logo_url, primary_color, subscription_status, created_at FROM merchants ORDER BY created_at DESC")
		if err != nil {
			utils.Error(c, 500, "failed to fetch merchants")
			return
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

		utils.Success(c, merchants)
	}
}

func AdminCreateMerchant(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			BusinessName  string `json:"business_name" binding:"required"`
			Slug          string `json:"slug"`
			AdminName     string `json:"admin_name"`
			AdminEmail    string `json:"admin_email" binding:"required,email"`
			AdminPassword string `json:"admin_password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body: "+err.Error())
			return
		}

		// Create merchant
		var merchantID string
		err := db.QueryRow(context.Background(),
			"INSERT INTO merchants (business_name, slug, subscription_status) VALUES ($1, $2, 'active') RETURNING id",
			req.BusinessName, req.Slug).Scan(&merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to create merchant")
			return
		}

		// Hash password
		hash, err := utils.HashPassword(req.AdminPassword)
		if err != nil {
			utils.Error(c, 500, "failed to process password")
			return
		}

		var userID string
		name := req.AdminName
		if name == "" {
			name = "Admin"
		}
		err = db.QueryRow(context.Background(),
			"INSERT INTO users (name, email, password_hash, role, merchant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			name, req.AdminEmail, hash, "MERCHANT_ADMIN", merchantID).Scan(&userID)
		if err != nil {
			utils.Error(c, 500, "failed to create admin user")
			return
		}

		c.JSON(201, gin.H{"success": true, "merchant_id": merchantID, "admin_id": userID})
	}
}

func AdminUpdateMerchant(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			BusinessName      string `json:"business_name"`
			SubscriptionStatus string `json:"subscription_status"`
			PrimaryColor      string `json:"primary_color"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
		}

		_, err := db.Exec(context.Background(),
			"UPDATE merchants SET business_name = COALESCE(NULLIF($1, ''), business_name), subscription_status = COALESCE(NULLIF($2, ''), subscription_status), primary_color = COALESCE(NULLIF($3, ''), primary_color) WHERE id = $4",
			req.BusinessName, req.SubscriptionStatus, req.PrimaryColor, id)
		if err != nil {
			utils.Error(c, 500, "failed to update merchant")
			return
		}

		utils.Message(c, "merchant updated")
	}
}

func AdminDeleteMerchant(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(context.Background(), "DELETE FROM merchants WHERE id = $1", id)
		if err != nil {
			utils.Error(c, 500, "failed to delete merchant")
			return
		}

		utils.Message(c, "merchant deleted")
	}
}

func AdminListAllSessions(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(),
			`SELECT s.id, s.merchant_id, s.status, s.payment_status, s.email, s.created_at, m.business_name 
			 FROM photobooth_sessions s
			 JOIN merchants m ON m.id = s.merchant_id
			 ORDER BY s.created_at DESC LIMIT 100`)
		if err != nil {
			utils.Error(c, 500, "failed to fetch sessions")
			return
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

		utils.Success(c, sessions)
	}
}

func AdminGetPlatformAnalytics(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var merchants, activeDevices, totalSessions, photos, prints, emails int
		var revenue float64

		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM merchants").Scan(&merchants)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM devices WHERE status = 'online'").Scan(&activeDevices)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM photobooth_sessions").Scan(&totalSessions)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM photos").Scan(&photos)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM print_jobs").Scan(&prints)
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM email_deliveries").Scan(&emails)
		db.QueryRow(context.Background(), "SELECT COALESCE(SUM(amount), 0) FROM payments").Scan(&revenue)

		utils.Success(c, map[string]interface{}{
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
