package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{Success: true, Data: data})
}

func SuccessWithMeta(c *gin.Context, data, meta interface{}) {
	c.JSON(200, APIResponse{Success: true, Data: data, Meta: meta})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{Success: false, Error: message})
}

func Message(c *gin.Context, message string) {
	c.JSON(200, APIResponse{Success: true, Message: message})
}
