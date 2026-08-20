package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"klikku/internal/utils"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.Error(c, fiber.StatusUnauthorized, "missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return utils.Error(c, fiber.StatusUnauthorized, "invalid authorization format")
	}

	cfg := utils.GetConfig(c)
	claims, err := utils.ValidateToken(tokenString, cfg)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "invalid or expired token")
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("merchant_id", claims.MerchantID)
	c.Locals("role", claims.Role)
	c.Locals("email", claims.Email)

	return c.Next()
}

func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}
		return utils.Error(c, fiber.StatusForbidden, "insufficient permissions")
	}
}

func SuperAdminMiddleware(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "SUPER_ADMIN" {
		return utils.Error(c, fiber.StatusForbidden, "super admin access required")
	}
	return c.Next()
}
