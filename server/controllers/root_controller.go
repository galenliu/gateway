package controllers

import (
	"gateway/pkg/log"
	"github.com/gofiber/fiber/v2"
)

func rootHandler(c *fiber.Ctx) error {
	log.Info("Root Handler Path:%s", c.Path())
	return c.Next()
}
