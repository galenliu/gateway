package controllers

import (
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type UserController struct {
	Users *models.Users
}

func NewUsersController() *UserController {
	uc := &UserController{}
	uc.Users = models.NewUsers()
	return uc
}

func (u *UserController) getCount(c *fiber.Ctx) error {
	count := u.Users.GetUsersCount()
	return c.SendString(strconv.Itoa(count))
}

func (u *UserController) createUser(c *fiber.Ctx) error {
	email := strings.ToLower(json.Get(c.Body(), "email").ToString())
	pw := json.Get(c.Body(), "password").ToString()
	name := json.Get(c.Body(), "password").ToString()

	if email == "" && pw == "" {
		return c.Status(fiber.StatusBadRequest).SendString("User requires email and password.")
	}
	exit := u.Users.GetUser(email)
	if exit != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User already exists.")
	}
	err, jwt := u.Users.CreateUser(email, pw, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(jwt)
}
