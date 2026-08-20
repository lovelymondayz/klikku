package utils

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Data: data})
}

func SuccessWithMeta(c *fiber.Ctx, data, meta interface{}) error {
	return c.JSON(APIResponse{Success: true, Data: data, Meta: meta})
}

func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(APIResponse{Success: false, Error: message})
}

func Message(c *fiber.Ctx, message string) error {
	return c.JSON(APIResponse{Success: true, Message: message})
}
