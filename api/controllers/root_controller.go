package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func staticHandler(c *fiber.Ctx) error {
	return nil
}

func rootHandler(c *fiber.Ctx) error {
	if c.Path() == "/" && c.Accepts("html") != "" {
		return staticHandler(c)
	}

	if c.Protocol() == "https" {
		c.Response().Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}
	c.Response().Header.Set("Content-Security-Policy", "")
	c.Response().Header.Set("Vary", "Accept")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")
	c.Response().Header.Set("Access-Control-Allow-Headers",
		"Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Response().Header.Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
	return c.Next()
}
