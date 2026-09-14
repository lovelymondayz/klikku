package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/utils"
)

func Register(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name         string `json:"name" binding:"required"`
			Email        string `json:"email" binding:"required,email"`
			Password     string `json:"password" binding:"required,min=6"`
			BusinessName string `json:"business_name"`
			Slug         string `json:"slug"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body: "+err.Error())
			return
		}

		// Check if email exists
		var existingID string
		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
		if err == nil {
			utils.Error(c, 409, "email already registered")
			return
		}

		// Hash password
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.Error(c, 500, "failed to process password")
			return
		}

		// Create merchant
		var merchantID string
		slug := req.Slug
		if slug == "" {
			slug = req.Email
		}
		err = db.QueryRow(context.Background(),
			"INSERT INTO merchants (business_name, slug, subscription_status) VALUES ($1, $2, 'active') RETURNING id",
			req.BusinessName, slug).Scan(&merchantID)
		if err != nil {
			utils.Error(c, 500, "failed to create merchant")
			return
		}

		// Create user
		var userID string
		err = db.QueryRow(context.Background(),
			"INSERT INTO users (name, email, password_hash, role, merchant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			req.Name, req.Email, hash, "MERCHANT_ADMIN", merchantID).Scan(&userID)
		if err != nil {
			utils.Error(c, 500, "failed to create user")
			return
		}

		utils.Success(c, gin.H{
			"user_id":     userID,
			"merchant_id": merchantID,
			"email":       req.Email,
			"role":        "MERCHANT_ADMIN",
		})
	}
}

func Login(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request body: "+err.Error())
			return
		}

		var user struct {
			ID           string
			Name         string
			Email        string
			PasswordHash string
			Role         string
			MerchantID   string
		}
		err := db.QueryRow(context.Background(),
			"SELECT id, name, email, password_hash, role, merchant_id FROM users WHERE email = $1",
			req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.MerchantID)
		if err != nil {
			utils.Error(c, 401, "invalid credentials")
			return
		}

		if !utils.CheckPassword(user.PasswordHash, req.Password) {
			utils.Error(c, 401, "invalid credentials")
			return
		}

		utils.Success(c, gin.H{
			"user_id":     user.ID,
			"merchant_id": user.MerchantID,
			"email":       user.Email,
			"name":        user.Name,
			"role":        user.Role,
		})
	}
}

func RefreshToken(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, 400, "invalid request")
			return
		}

		utils.Message(c, "not yet implemented")
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Message(c, "logged out successfully")
	}
}

func GetCurrentUser(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")
		role, _ := c.Get("role")
		merchantID, _ := c.Get("merchant_id")

		utils.Success(c, gin.H{
			"user_id":     userID,
			"email":       email,
			"role":        role,
			"merchant_id": merchantID,
		})
	}
}
