package controllers

import (
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type LoginController struct {
	store UsersStore
}

func NewLoginController(store UsersStore) *LoginController {
	c := LoginController{}
	c.store = store
	return &c
}

func (c LoginController) handleLogin(ctx *fiber.Ctx) error {
	email := strings.ToLower(json.Get(ctx.Body(), "email").ToString())
	password := strings.ToLower(json.Get(ctx.Body(), "email").ToString())

	if email == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	user := c.store.GetUserByEmail(email)
	if user == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), pwHash)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	return nil
}
