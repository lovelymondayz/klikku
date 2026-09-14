package middleware

import (
	"github.com/gin-gonic/gin"
	"klikku/internal/utils"
)

func TenantMiddleware(c *gin.Context) {
	role, _ := c.Get("role")
	merchantID, _ := c.Get("merchant_id")

	if role == "SUPER_ADMIN" {
		c.Next()
		return
	}

	if merchantID == nil || merchantID.(string) == "" {
		utils.Error(c, 403, "no merchant context")
		c.Abort()
		return
	}

	c.Set("tenant_id", merchantID.(string))
	c.Next()
}
