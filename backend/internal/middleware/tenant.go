package middleware

import (
	"github.com/gofiber/fiber/v2"
	"klikku/internal/utils"
)

func TenantMiddleware(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	merchantID := c.Locals("merchant_id").(string)

	// Super Admin bypasses tenant isolation
	if role == "SUPER_ADMIN" {
		return c.Next()
	}

	if merchantID == "" {
		return utils.Error(c, fiber.StatusForbidden, "no merchant context")
	}

	c.Locals("tenant_id", merchantID)
	return c.Next()
}
