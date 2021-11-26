package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/xiam/to"
	"sync"
)

const (
	ActionCompleted = "completed"
	ActionPending   = "pending"
	ActionCreated   = "created"
	ActionError     = "error"
	ActionDeleted   = "deleted"

	ActionPair   = "pair"
	ActionUnpair = "unpair"
)

type ActionsManager interface {
	AddNewThings(timeOut int) error
	RequestAction(thingId, actionName string, input map[string]interface{}) error
	RemoveThing(id string) error
}

type ActionsModel struct {
	logger    logging.Logger
	Actions   sync.Map
	container *container.ThingsContainer
	manager   ActionsManager
}

func NewActionsModel(m ActionsManager, container *container.ThingsContainer, log logging.Logger) *ActionsModel {
	return &ActionsModel{
		logger:    log,
		container: container,
		manager:   m,
	}
}

func (m *ActionsModel) Add(a *Action) error {
	m.Actions.Store(a.Id, a)
	a.updateStatus(ActionPending)
	switch a.Name {
	case ActionPair:
		timeout := to.Int(a.Input["timeout"])
		err := m.manager.AddNewThings(timeout)
		if err != nil {
			a.Error = err
			a.updateStatus(ActionError)
			m.logger.Infof("things was not added.err:%s", err.Error())
			return err
		}
		a.updateStatus(ActionCompleted)
	case ActionUnpair:
		thingId, _ := a.GetInput()["id"].(string)
		thing := m.container.GetThing(thingId)
		if thing == nil {
			err := fmt.Errorf("requset parems err")
			a.SetErr(err)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		err := m.manager.RemoveThing(thingId)
		err = m.container.RemoveThing(thingId)
		if err != nil {

		}
	default:

	}
	return nil
}

func (m *ActionsModel) OnActionStatus(model *Action) {

}
