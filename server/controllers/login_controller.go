package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
)

type LoginController struct {
	users  *models.Users
	logger logging.Logger
}

func NewLoginController(users *models.Users, logger logging.Logger) *LoginController {
	c := LoginController{}
	c.logger = logger
	c.users = users
	return &c
}

func (c LoginController) handleLogin(ctx *fiber.Ctx) error {
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")

	if email == "" || password == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	user := c.users.GetUser(email)
	if user == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	if !user.ComparePassword(password) {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	return nil
}
