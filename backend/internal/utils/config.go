package utils

import (
	"github.com/gofiber/fiber/v2"
	"klikku/internal/config"
)

// Store config in fiber locals during app setup
func SetConfig(c *fiber.Ctx, cfg *config.Config) {
	c.Locals("config", cfg)
}

func GetConfig(c *fiber.Ctx) *config.Config {
	return c.Locals("config").(*config.Config)
}
