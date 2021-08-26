package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type JWT interface {
	IssueToken(user int64) string
}

type UserController struct {
	model  *models.Users
	logger logging.Logger
	jwt    JWT
}

func NewUsersController(m *models.Users, jwt JWT, log logging.Logger) *UserController {
	uc := &UserController{}
	uc.jwt = jwt
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
	id, err := u.model.CreateUser(email, pw, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	jwt := u.jwt.IssueToken(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": jwt})

}
