package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/utils"
)

type CreateCampaignRequest struct {
	Name            string                 `json:"name" binding:"required"`
	Description     string                 `json:"description"`
	StartDate       time.Time              `json:"start_date" binding:"required"`
	EndDate         time.Time              `json:"end_date" binding:"required"`
	Status          string                 `json:"status"`
	PromotionConfig map[string]interface{} `json:"promotion_config"`
}

type CreateTemplateRequest struct {
	Name         string                 `json:"name" binding:"required"`
	CampaignID   string                 `json:"campaign_id"`
	LayoutConfig map[string]interface{} `json:"layout_config"`
	OverlayURL   string                 `json:"overlay_url"`
	OutputWidth  int                    `json:"output_width"`
	OutputHeight int                    `json:"output_height"`
	PhotoCount   int                    `json:"photo_count"`
	Price        float64                `json:"price"`
	Active       bool                   `json:"active"`
}

func ListCampaigns(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, merchant_id, name, description, start_date, end_date, status, promotion_config, created_at FROM campaigns WHERE merchant_id = $1 ORDER BY created_at DESC",
			merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch campaigns")
			return
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

		utils.Success(c, campaigns)
	}
}

func CreateCampaign(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		var req CreateCampaignRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body: "+err.Error())
			return
		}

		promotionConfig := req.PromotionConfig
		if promotionConfig == nil {
			promotionConfig = make(map[string]interface{})
		}

		var id string
		err := db.QueryRow(context.Background(),
			`INSERT INTO campaigns (merchant_id, name, description, start_date, end_date, status, promotion_config) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			merchantID, req.Name, req.Description, req.StartDate, req.EndDate, req.Status, promotionConfig).Scan(&id)
		if err != nil {
			utils.Error(c, 500, "failed to create campaign")
			return
		}

		c.JSON(201, gin.H{"success": true, "id": id})
	}
}

func GetCampaign(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var mID, name, description, status string
		var startDate, endDate, createdAt time.Time
		var promotionConfig map[string]interface{}

		err := db.QueryRow(context.Background(),
			"SELECT id, merchant_id, name, description, start_date, end_date, status, promotion_config, created_at FROM campaigns WHERE id = $1",
			id).Scan(&id, &mID, &name, &description, &startDate, &endDate, &status, &promotionConfig, &createdAt)
		if err != nil {
			utils.Error(c, 404, "campaign not found")
			return
		}

		utils.Success(c, map[string]interface{}{
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

func UpdateCampaign(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Name            string                 `json:"name"`
			Description     string                 `json:"description"`
			Status          string                 `json:"status"`
			PromotionConfig map[string]interface{} `json:"promotion_config"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
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
			utils.Error(c, 400, "no fields to update")
			return
		}

		query += strings.Join(updates, ", ") + " WHERE id = $"+string(rune('0'+argIdx))
		args = append(args, id)

		_, err := db.Exec(context.Background(), query, args...)
		if err != nil {
			utils.Error(c, 500, "failed to update campaign")
			return
		}

		utils.Message(c, "campaign updated")
	}
}

func DeleteCampaign(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(context.Background(), "DELETE FROM campaigns WHERE id = $1", id)
		if err != nil {
			utils.Error(c, 500, "failed to delete campaign")
			return
		}

		utils.Message(c, "campaign deleted")
	}
}

func ListTemplates(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		rows, err := db.Query(context.Background(),
			"SELECT id, merchant_id, campaign_id, name, preview_url, layout_config, overlay_url, output_width, output_height, photo_count, price, active, created_at FROM templates WHERE merchant_id = $1 ORDER BY created_at DESC",
			merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to fetch templates")
			return
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

		utils.Success(c, templates)
	}
}

func CreateTemplate(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID := getMerchantID(c)
		if merchantID == "" {
			utils.Error(c, 400, "merchant_id required")
			return
		}

		var req CreateTemplateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body: "+err.Error())
			return
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
			utils.Error(c, 500, "failed to create template")
			return
		}

		c.JSON(201, gin.H{"success": true, "id": id})
	}
}

func GetTemplate(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var mID, campaignID, name, previewURL, overlayURL string
		var layoutConfig map[string]interface{}
		var outputWidth, outputHeight, photoCount int
		var price float64
		var active bool
		var createdAt time.Time

		err := db.QueryRow(context.Background(),
			"SELECT id, merchant_id, campaign_id, name, preview_url, layout_config, overlay_url, output_width, output_height, photo_count, price, active, created_at FROM templates WHERE id = $1",
			id).Scan(&id, &mID, &campaignID, &name, &previewURL, &layoutConfig, &overlayURL, &outputWidth, &outputHeight, &photoCount, &price, &active, &createdAt)
		if err != nil {
			utils.Error(c, 404, "template not found")
			return
		}

		utils.Success(c, map[string]interface{}{
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

func UpdateTemplate(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Name         string                 `json:"name"`
			PhotoCount   int                    `json:"photo_count"`
			Price        float64                `json:"price"`
			LayoutConfig map[string]interface{} `json:"layout_config"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body")
			return
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
			utils.Error(c, 400, "no fields to update")
			return
		}

		query += strings.Join(updates, ", ") + " WHERE id = $"+string(rune('0'+argIdx))
		args = append(args, id)

		_, err := db.Exec(context.Background(), query, args...)
		if err != nil {
			utils.Error(c, 500, "failed to update template")
			return
		}

		utils.Message(c, "template updated")
	}
}

func DeleteTemplate(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(context.Background(), "DELETE FROM templates WHERE id = $1", id)
		if err != nil {
			utils.Error(c, 500, "failed to delete template")
			return
		}

		utils.Message(c, "template deleted")
	}
}
