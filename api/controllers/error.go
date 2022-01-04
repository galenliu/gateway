package controllers

import "github.com/gofiber/fiber/v2"

func errBadRequest(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": msg,
	})
}
