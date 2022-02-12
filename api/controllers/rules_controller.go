package controllers

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rules_engine"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/fiber/v2"
)

type RulesController struct {
	engine *rules_engine.Engine
	logger logging.Logger
}

func NewRulesController(db rules_engine.RuleDB, container things.Container) *RulesController {
	c := &RulesController{}
	c.engine = rules_engine.NewEngine(db, container)
	return c
}

func (r *RulesController) handleGetRules(ctx *fiber.Ctx) error {
	return nil
}

func (r *RulesController) handleCreateRule(c *fiber.Ctx) error {
	fmt.Printf("Post /rule:\t\n %s\n", c.Body())
	rule, err := r.engine.CreateRule(c.Body())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(util.JsonIndent(rule))
}
