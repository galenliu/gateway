package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func NotFoundError(info string, args ...string) error {
	return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf(info, args))
}

func BadRequestError(info string, args ...string) error {
	return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(info, args))
}
