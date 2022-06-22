package controllers

import (
	"github.com/galenliu/gateway/api/models"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type JWT interface {
	IssueToken(user int64) (string, error)
}

type UserController struct {
	model  *models.Users
	logger log.Logger
	jwt    JWT
}

func NewUsersController(m *models.Users, jwt JWT) *UserController {
	uc := &UserController{}
	uc.jwt = jwt
	uc.model = m

	return uc
}

func (u *UserController) getCount(c *fiber.Ctx) error {
	count := u.model.GetUsersCount()
	return c.SendString(strconv.Itoa(count))
}

func (u *UserController) createUser(c *fiber.Ctx) error {
	email := c.FormValue("email")
	pw := c.FormValue("password")
	name := c.FormValue("name")

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
	jwt, err := u.jwt.IssueToken(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": jwt})

}
