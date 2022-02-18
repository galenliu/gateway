package controllers

import "github.com/gofiber/fiber/v2"

type GroupsController struct {
}

func NewGroupsController() *GroupsController {
	return &GroupsController{}
}

func (c GroupsController) handleGetGroups(ctx *fiber.Ctx) error {
	return ctx.JSON(struct{}{})
}

func (c GroupsController) handleGetGroup(ctx *fiber.Ctx) error {
	return ctx.JSON(struct{}{})
}

func (c GroupsController) handlerCreateGroup(ctx *fiber.Ctx) error {
	return ctx.JSON(struct{}{})
}

func (c GroupsController) handlerDeleteGroup(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(id)
}

func (c GroupsController) handlerUpdateGroup(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(id)
}
