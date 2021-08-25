package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type UserController struct {
	model  *models.Users
	logger logging.Logger
}

func NewUsersController(m *models.Users, log logging.Logger) *UserController {
	uc := &UserController{}
	uc.model = m
	uc.logger = log
	return uc
}

func (u *UserController) getCount(c *fiber.Ctx) error {
	count := u.model.GetUsersCount()
	return c.SendString(strconv.Itoa(count))
}

func (u *UserController) createUser(c *fiber.Ctx) error {
	email := strings.ToLower(json.Get(c.Body(), "email").ToString())
	pw := json.Get(c.Body(), "password").ToString()
	name := json.Get(c.Body(), "name").ToString()

	if email == "" && pw == "" {
		return c.Status(fiber.StatusBadRequest).SendString("User requires email and password.")
	}
	exit := u.model.GetUser(email)
	if exit != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User already exists.")
	}
	jwt, err := u.model.CreateUser(email, pw, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(jwt)
}
