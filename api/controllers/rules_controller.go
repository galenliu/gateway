package controllers

import (
	"encoding/json"
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/rules_engine"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type RulesController struct {
	engine *rules_engine.Engine
	logger log.Logger
}

func NewRulesController(db rules_engine.RuleDB, container things.Container) *RulesController {
	c := &RulesController{}
	c.engine = rules_engine.NewEngine(db, container)
	return c
}

func (r *RulesController) handleGetRules(ctx *fiber.Ctx) error {
	rules := make([]*rules_engine.Rule, 0)
	rs := r.engine.GetRules()
	for _, r := range rs {
		rules = append(rules, r)
	}
	return ctx.JSON(rules)
}

func (r *RulesController) handleCreateRule(c *fiber.Ctx) error {
	fmt.Printf("Post /rule:\t\n %s\n", c.Body())
	var desc rules_engine.RuleDescription
	err := json.Unmarshal(c.Body(), &desc)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	rule, err := r.engine.AddRule(desc)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(util.JsonIndent(rule))
}

func (r *RulesController) handleGetRule(ctx *fiber.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	rule := r.engine.GetRule(int64(id))
	return ctx.JSON(rule)
}

func (r *RulesController) handlerDeleteRule(ctx *fiber.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	err = r.engine.DeleteRule(int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (r *RulesController) handlerUpdateRule(ctx *fiber.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}
	var desc rules_engine.RuleDescription
	err = json.Unmarshal(ctx.Body(), &desc)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err = r.engine.UpdateRule(int64(id), desc)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.SendStatus(fiber.StatusOK)
}
