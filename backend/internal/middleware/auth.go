package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"klikku/internal/utils"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Error(c, 401, "missing authorization header")
		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		utils.Error(c, 401, "invalid authorization format")
		c.Abort()
		return
	}

	cfg := utils.GetConfig(c)
	claims, err := utils.ValidateToken(tokenString, cfg)
	if err != nil {
		utils.Error(c, 401, "invalid or expired token")
		c.Abort()
		return
	}

	c.Set("user_id", claims.UserID)
	c.Set("merchant_id", claims.MerchantID)
	c.Set("role", claims.Role)
	c.Set("email", claims.Email)
	c.Next()
}

func SuperAdminMiddleware(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role.(string) != "SUPER_ADMIN" {
		utils.Error(c, 403, "super admin access required")
		c.Abort()
		return
	}
	c.Next()
}
