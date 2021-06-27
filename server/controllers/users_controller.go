package controllers

import (
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type usesModel interface {
	GetUsersCount() int
}

type userController struct {
	model usesModel
}

func NewUsersController(model usesModel) *userController {
	uc := &userController{}
	uc.model = model
	return uc
}

func (u *userController) getCount(c *fiber.Ctx) error {
	count := u.model.GetUsersCount()
	return c.SendString(strconv.Itoa(count))
}

func (u *userController) createUser(c *fiber.Ctx) error {
	email := strings.ToLower(json.Get(c.Body(), "email").ToString())
	pw := json.Get(c.Body(), "password").ToString()
	name := json.Get(c.Body(), "password").ToString()

	if email == "" && pw == "" {
		return c.Status(fiber.StatusBadRequest).SendString("User requires email and password.")
	}
	exit := u.model.GetUser(email)
	if exit != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User already exists.")
	}
	err, jwt := u.Users.CreateUser(email, pw, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(jwt)
}
