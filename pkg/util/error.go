package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func NotFoundError(info string, args ...string) error {
	return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf(info, args))
}

func BadRequestError(info string, args ...string) error {
	return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(info, args))
}

func SuccessOrExit(err error) {
	if err != nil {
		log.Println("exit()")
		log.Fatal(err.Error())
	}
}

func LogError(err error, name, msg string) {
	if err != nil {
		log.Print(name + ":")
		log.Printf(msg, err.Error())
	}
}
