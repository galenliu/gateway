package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func rootHandler(c *fiber.Ctx) error {
	return c.Next()
}
