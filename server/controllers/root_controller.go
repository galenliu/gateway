package controllers

import (
	"gateway/pkg/log"
	"github.com/gofiber/fiber/v2"
)

func rootHandler(c *fiber.Ctx) error {
	log.Info( "Content-Type: %s",c.Get("Content-Type"))
	return c.Next()
}
