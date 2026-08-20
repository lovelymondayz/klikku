package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/utils"
)

// Request/Response types
type CreateCampaignRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	StartDate       time.Time              `json:"start_date"`
	EndDate         time.Time              `json:"end_date"`
	Status          string                 `json:"status"`
	PromotionConfig map[string]interface{} `json:"promotion_config"`
}

type UpdateCampaignRequest struct {
	Name            string                 `json:"name,omitempty"`
	Description     string                 `json:"description,omitempty"`
	StartDate       *time.Time             `json:"start_date,omitempty"`
	EndDate         *time.Time             `json:"end_date,omitempty"`
	Status          string                 `json:"status,omitempty"`
	PromotionConfig map[string]interface{} `json:"promotion_config,omitempty"`
}

type CreateTemplateRequest struct {
	Name         string                 `json:"name"`
	CampaignID   string                 `json:"campaign_id,omitempty"`
	LayoutConfig map[string]interface{} `json:"layout_config"`
	OverlayURL   string                 `json:"overlay_url,omitempty"`
	OutputWidth  int                    `json:"output_width"`
	OutputHeight int                    `json:"output_height"`
	PhotoCount   int                    `json:"photo_count"`
	Price        float64                `json:"price"`
	Active       bool                   `json:"active"`
}

type UpdateTemplateRequest struct {
	Name         string                 `json:"name,omitempty"`
	CampaignID   string                 `json:"campaign_id,omitempty"`
	LayoutConfig map[string]interface{} `json:"layout_config,omitempty"`
	OverlayURL   string                 `json:"overlay_url,omitempty"`
	OutputWidth  int                    `json:"output_width,omitempty"`
	OutputHeight int                    `json:"output_height,omitempty"`
	PhotoCount   int                    `json:"photo_count,omitempty"`
	Price        float64                `json:"price,omitempty"`
	Active       *bool                  `json:"active,omitempty"`
}

func getMerchantID(c *fiber.Ctx) string {
	role := c.Locals("role").(string)
	if role == "SUPER_ADMIN" {
		return c.Query("merchant_id", "")
	}
	return c.Locals("merchant_id").(string)
}

// Campaigns
func ListCampaigns(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, merchant_id, name, description, start_date, end_date, status, promotion_config, created_at FROM campaigns WHERE merchant_id = $1 ORDER BY created_at DESC",
			merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch campaigns")
		}
		defer rows.Close()

		var campaigns []map[string]interface{}
		for rows.Next() {
			var id, mID, name, description, status string
			var startDate, endDate, createdAt time.Time
			var promotionConfig map[string]interface{}
			err := rows.Scan(&id, &mID, &name, &description, &startDate, &endDate, &status, &promotionConfig, &createdAt)
			if err != nil {
				continue
			}
			campaigns = append(campaigns, map[string]interface{}{
				"id":               id,
				"merchant_id":      mID,
				"name":             name,
				"description":      description,
				"start_date":       startDate,
				"end_date":         endDate,
				"status":           status,
				"promotion_config": promotionConfig,
				"created_at":       createdAt,
			})
		}

		return utils.Success(c, campaigns)
	}
}

func CreateCampaign(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		var req CreateCampaignRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		if req.Name == "" {
			return utils.Error(c, fiber.StatusBadRequest, "campaign name required")
		}

		status := req.Status
		if status == "" {
			status = "draft"
		}

		promotionConfig := req.PromotionConfig
		if promotionConfig == nil {
			promotionConfig = make(map[string]interface{})
		}

		var id string
		err := db.QueryRow(context.Background(),
			`INSERT INTO campaigns (merchant_id, name, description, start_date, end_date, status, promotion_config) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			merchantID, req.Name, req.Description, req.StartDate, req.EndDate, status, promotionConfig).Scan(&id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create campaign")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "id": id})
	}
}

func GetCampaign(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid campaign ID")
		}

		var mID, name, description, status string
		var startDate, endDate, createdAt time.Time
		var promotionConfig map[string]interface{}

		err = db.QueryRow(context.Background(),
			"SELECT id, merchant_id, name, description, start_date, end_date, status, promotion_config, created_at FROM campaigns WHERE id = $1",
			id).Scan(&id, &mID, &name, &description, &startDate, &endDate, &status, &promotionConfig, &createdAt)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "campaign not found")
		}

		return utils.Success(c, map[string]interface{}{
			"id":               id,
			"merchant_id":      mID,
			"name":             name,
			"description":      description,
			"start_date":       startDate,
			"end_date":         endDate,
			"status":           status,
			"promotion_config": promotionConfig,
			"created_at":       createdAt,
		})
	}
}

func UpdateCampaign(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid campaign ID")
		}

		var req UpdateCampaignRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		query := "UPDATE campaigns SET "
		args := []interface{}{}
		argIdx := 1
		updates := []string{}

		if req.Name != "" {
			updates = append(updates, "name = $"+string(rune('0'+argIdx)))
			args = append(args, req.Name)
			argIdx++
		}
		if req.Description != "" {
			updates = append(updates, "description = $"+string(rune('0'+argIdx)))
			args = append(args, req.Description)
			argIdx++
		}
		if req.Status != "" {
			updates = append(updates, "status = $"+string(rune('0'+argIdx)))
			args = append(args, req.Status)
			argIdx++
		}

		if len(updates) == 0 {
			return utils.Error(c, fiber.StatusBadRequest, "no fields to update")
		}

		query += strings.Join(updates, ", ") + " WHERE id = $" + string(rune('0'+argIdx))
		args = append(args, id)

		_, err = db.Exec(context.Background(), query, args...)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to update campaign")
		}

		return utils.Message(c, "campaign updated")
	}
}

func DeleteCampaign(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid campaign ID")
		}

		_, err = db.Exec(context.Background(), "DELETE FROM campaigns WHERE id = $1", id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to delete campaign")
		}

		return utils.Message(c, "campaign deleted")
	}
}

// Templates
func ListTemplates(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, merchant_id, campaign_id, name, preview_url, layout_config, overlay_url, output_width, output_height, photo_count, price, active, created_at FROM templates WHERE merchant_id = $1 ORDER BY created_at DESC",
			merchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to fetch templates")
		}
		defer rows.Close()

		var templates []map[string]interface{}
		for rows.Next() {
			var id, mID, campaignID, name, previewURL, overlayURL string
			var layoutConfig map[string]interface{}
			var outputWidth, outputHeight, photoCount int
			var price float64
			var active bool
			var createdAt time.Time
			err := rows.Scan(&id, &mID, &campaignID, &name, &previewURL, &layoutConfig, &overlayURL, &outputWidth, &outputHeight, &photoCount, &price, &active, &createdAt)
			if err != nil {
				continue
			}
			templates = append(templates, map[string]interface{}{
				"id":             id,
				"merchant_id":    mID,
				"campaign_id":    campaignID,
				"name":           name,
				"preview_url":    previewURL,
				"layout_config":  layoutConfig,
				"overlay_url":    overlayURL,
				"output_width":   outputWidth,
				"output_height":  outputHeight,
				"photo_count":    photoCount,
				"price":          price,
				"active":         active,
				"created_at":     createdAt,
			})
		}

		return utils.Success(c, templates)
	}
}

func CreateTemplate(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			return utils.Error(c, fiber.StatusBadRequest, "merchant_id required")
		}

		var req CreateTemplateRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		if req.Name == "" {
			return utils.Error(c, fiber.StatusBadRequest, "template name required")
		}

		layoutConfig := req.LayoutConfig
		if layoutConfig == nil {
			layoutConfig = make(map[string]interface{})
		}

		outputWidth := req.OutputWidth
		if outputWidth == 0 {
			outputWidth = 1200
		}
		outputHeight := req.OutputHeight
		if outputHeight == 0 {
			outputHeight = 1800
		}
		photoCount := req.PhotoCount
		if photoCount == 0 {
			photoCount = 4
		}

		var id string
		err := db.QueryRow(context.Background(),
			`INSERT INTO templates (merchant_id, campaign_id, name, layout_config, overlay_url, output_width, output_height, photo_count, price, active) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			merchantID, req.CampaignID, req.Name, layoutConfig, req.OverlayURL, outputWidth, outputHeight, photoCount, req.Price, req.Active).Scan(&id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create template")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "id": id})
	}
}

func GetTemplate(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid template ID")
		}

		var mID, campaignID, name, previewURL, overlayURL string
		var layoutConfig map[string]interface{}
		var outputWidth, outputHeight, photoCount int
		var price float64
		var active bool
		var createdAt time.Time

		err = db.QueryRow(context.Background(),
			"SELECT id, merchant_id, campaign_id, name, preview_url, layout_config, overlay_url, output_width, output_height, photo_count, price, active, created_at FROM templates WHERE id = $1",
			id).Scan(&id, &mID, &campaignID, &name, &previewURL, &layoutConfig, &overlayURL, &outputWidth, &outputHeight, &photoCount, &price, &active, &createdAt)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "template not found")
		}

		return utils.Success(c, map[string]interface{}{
			"id":             id,
			"merchant_id":    mID,
			"campaign_id":    campaignID,
			"name":           name,
			"preview_url":    previewURL,
			"layout_config":  layoutConfig,
			"overlay_url":    overlayURL,
			"output_width":   outputWidth,
			"output_height":  outputHeight,
			"photo_count":    photoCount,
			"price":          price,
			"active":         active,
			"created_at":     createdAt,
		})
	}
}

func UpdateTemplate(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid template ID")
		}

		var req UpdateTemplateRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		query := "UPDATE templates SET "
		args := []interface{}{}
		argIdx := 1
		updates := []string{}

		if req.Name != "" {
			updates = append(updates, "name = $"+string(rune('0'+argIdx)))
			args = append(args, req.Name)
			argIdx++
		}
		if req.PhotoCount > 0 {
			updates = append(updates, "photo_count = $"+string(rune('0'+argIdx)))
			args = append(args, req.PhotoCount)
			argIdx++
		}

		if len(updates) == 0 {
			return utils.Error(c, fiber.StatusBadRequest, "no fields to update")
		}

		query += strings.Join(updates, ", ") + " WHERE id = $" + string(rune('0'+argIdx))
		args = append(args, id)

		_, err = db.Exec(context.Background(), query, args...)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to update template")
		}

		return utils.Message(c, "template updated")
	}
}

func DeleteTemplate(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid template ID")
		}

		_, err = db.Exec(context.Background(), "DELETE FROM templates WHERE id = $1", id)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to delete template")
		}

		return utils.Message(c, "template deleted")
	}
}
