package controllers

import (
	"github.com/galenliu/gateway/config"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/gofiber/fiber/v2"
)

func rootHandler(c *fiber.Ctx) error {
	if config.IsVerbose() {
		log.Info("url:%s path: %s", c.BaseURL(), c.Path())
	}
	return c.Next()
}
