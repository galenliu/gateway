package controllers

import (
	"github.com/galenliu/gateway/api/models"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/gofiber/fiber/v2"
)

type Users interface {
	GetUser(email string) *models.User
}

type LoginController struct {
	users        *models.Users
	jsonwebtoken *models.Jsonwebtoken
	logger       log.Logger
}

func NewLoginController(users *models.Users, jsonwebtoken *models.Jsonwebtoken) *LoginController {
	c := LoginController{}
	c.users = users
	c.jsonwebtoken = jsonwebtoken

	c.users = users
	return &c
}

func (c *LoginController) handleLogin(ctx *fiber.Ctx) error {
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")
	if email == "" || password == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	user := c.users.GetUser(email)
	if user == nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("user not exist")
	}
	if !user.ComparePassword(password) {
		return ctx.Status(fiber.StatusUnauthorized).SendString("password invalid")
	}
	jwt, err := c.jsonwebtoken.IssueToken(user.ID)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": jwt})
}
