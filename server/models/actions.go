package models

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
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
	CancelAddNewThing()
	RequestAction(ctx context.Context, thingId, actionName string, input map[string]interface{}) error
	RemoveThing(id string) error
	RemoveAction(thingId string, actionId string, actionName string) error
	CancelRemoveThing(id string)
}

type ActionsModel struct {
	logger    logging.Logger
	actions   sync.Map
	container *container.ThingsContainer
	manager   ActionsManager
	bus       *bus.Bus
}

func NewActionsModel(m ActionsManager, container *container.ThingsContainer, bus *bus.Bus, log logging.Logger) *ActionsModel {
	return &ActionsModel{
		bus:       bus,
		logger:    log,
		container: container,
		manager:   m,
	}
}

func (m *ActionsModel) Add(a *Action) error {
	m.actions.Store(a.GetId(), a)
	m.onActionStatus(a)
	if a.GetThingId() != "" {
		thing := m.container.GetThing(a.GetThingId())
		success := thing.AddAction(a.GetName())
		if !success {
			m.actions.Delete(a.GetId())
			return fmt.Errorf("invalid thing action name: %s", a.GetName())
		}
	}
	a.updateStatus(ActionPending)
	switch a.GetName() {
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
		return nil
	case ActionUnpair:
		thingId, _ := a.GetInput()["id"].(string)
		thing := m.container.GetThing(thingId)
		if thing == nil {
			err := fmt.Errorf("requset parems err")
			a.SetErr(err)
			m.logger.Error(err)
			return err
		}
		err := m.manager.RemoveThing(thingId)
		if err != nil {
			err := fmt.Errorf("addon unpair thing: %s failed", thingId)
			a.SetErr(err)
			m.logger.Error(err)
			return err
		}
		err = m.container.RemoveThing(thingId)
		if err != nil {
			err := fmt.Errorf("unpair of thing: %s failed", thingId)
			a.SetErr(err)
			m.logger.Error(err)
			return err
		}
	default:
		m.actions.Delete(a.GetId())
		return fmt.Errorf("invalid action name: %s", a.GetName())
	}
	return nil
}

func (m *ActionsModel) onActionStatus(a *Action) {
	m.bus.Pub(topic.ThingActionStatus, a)
}

func (m *ActionsModel) updateStatus(ad *addon.ActionDescription) {
	a, ok := m.actions.Load(ad.Id)
	if !ok {
		return
	}
	action, ok := a.(*Action)
	action.update(ad)
}

func (m *ActionsModel) Remove(id string) error {
	action, _ := m.actions.Load(id)
	a, ok := action.(*Action)
	if action == nil || !ok {
		return fmt.Errorf("invalid action id: %s", id)
	}
	defer func() {
		a.updateStatus(ActionDeleted)
		m.actions.Delete(id)
	}()
	if a.GetStatus() == ActionPending {
		if a.GetThingId() != "" {
			thing := m.container.GetThing(a.GetThingId())
			if thing != nil {
				if thing.RemoveAction(a.GetName()) {
					return fmt.Errorf("invalid action name %s", a.GetName())
				}
			} else {
				return fmt.Errorf("error removing thing action: %s", a.GetName())
			}
		} else {
			return fmt.Errorf("error removing thing action: %s", a.GetName())
		}
	}
	switch a.GetName() {
	case ActionPair:
		m.manager.CancelAddNewThing()
		break
	case ActionUnpair:
		id = a.GetInput()["id"].(string)
		if id == "" {
			return fmt.Errorf("unpair id invalid ")
		}
		m.manager.CancelRemoveThing(id)
	default:
		return fmt.Errorf("invaild action name: %s", a.GetName())
	}
	return nil
}
