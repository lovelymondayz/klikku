package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"klikku/internal/config"
	"klikku/internal/models"
	"klikku/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	BusinessName string `json:"business_name"`
	Slug         string `json:"slug"`
}

func Register(db *pgxpool.Pool, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		if req.Email == "" || req.Password == "" || req.Name == "" {
			return utils.Error(c, fiber.StatusBadRequest, "missing required fields")
		}

		// Check if email exists
		var existingID string
		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
		if err == nil {
			return utils.Error(c, fiber.StatusConflict, "email already registered")
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to process password")
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
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create merchant")
		}

		// Create user
		var userID string
		err = db.QueryRow(context.Background(),
			"INSERT INTO users (name, email, password_hash, role, merchant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			req.Name, req.Email, string(hash), models.RoleMerchantAdmin, merchantID).Scan(&userID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "failed to create user")
		}

		// Generate tokens
		accessToken, _ := utils.GenerateToken(userID, merchantID, models.RoleMerchantAdmin, req.Email, cfg)
		refreshToken, _ := utils.GenerateRefreshToken(userID, cfg)

		return utils.Success(c, fiber.Map{
			"user_id":       userID,
			"merchant_id":   merchantID,
			"email":         req.Email,
			"role":          models.RoleMerchantAdmin,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func Login(db *pgxpool.Pool, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
		}

		var user models.User
		err := db.QueryRow(context.Background(),
			"SELECT id, name, email, password_hash, role, merchant_id FROM users WHERE email = $1",
			req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.MerchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "invalid credentials")
		}

		if !utils.CheckPassword(user.PasswordHash, req.Password) {
			return utils.Error(c, fiber.StatusUnauthorized, "invalid credentials")
		}

		accessToken, _ := utils.GenerateToken(user.ID, user.MerchantID, user.Role, user.Email, cfg)
		refreshToken, _ := utils.GenerateRefreshToken(user.ID, cfg)

		return utils.Success(c, fiber.Map{
			"user_id":       user.ID,
			"merchant_id":   user.MerchantID,
			"email":         user.Email,
			"name":          user.Name,
			"role":          user.Role,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken(db *pgxpool.Pool, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request")
		}

		claims, err := utils.ValidateToken(req.RefreshToken, cfg)
		if err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "invalid refresh token")
		}

		var user models.User
		err = db.QueryRow(context.Background(),
			"SELECT id, name, email, role, merchant_id FROM users WHERE id = $1",
			claims.UserID).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.MerchantID)
		if err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "user not found")
		}

		accessToken, _ := utils.GenerateToken(user.ID, user.MerchantID, user.Role, user.Email, cfg)
		newRefreshToken, _ := utils.GenerateRefreshToken(user.ID, cfg)

		return utils.Success(c, fiber.Map{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
		})
	}
}

func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.Message(c, "logged out successfully")
	}
}

func GetCurrentUser(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)

		var user models.User
		err := db.QueryRow(context.Background(),
			"SELECT id, name, email, role, merchant_id, created_at FROM users WHERE id = $1",
			userID).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.MerchantID, &user.CreatedAt)
		if err != nil {
			return utils.Error(c, fiber.StatusNotFound, "user not found")
		}

		return utils.Success(c, user)
	}
}
