package utils

import (
	"github.com/gin-gonic/gin"
	"klikku/internal/config"
)

func SetConfig(c *gin.Context, cfg *config.Config) {
	c.Set("config", cfg)
}

func GetConfig(c *gin.Context) *config.Config {
	return c.MustGet("config").(*config.Config)
}
