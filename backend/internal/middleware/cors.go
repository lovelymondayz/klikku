package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"klikku/internal/config"
)

func CORS(cfg *config.Config) cors.Config {
	return cors.Config{
		AllowOrigins:     cfg.FrontendURL + ",http://localhost:3009,http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Device-Token",
		AllowCredentials: true,
	}
}
