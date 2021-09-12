package models

import (
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin"
	"github.com/xiam/to"
	"sync"
)

type ActionsManager interface {
	AddNewThings(timeOut int) error
}

type ActionsBus struct {
	bus.Controller
}

type ActionsModel struct {
	logger  logging.Logger
	Actions sync.Map
	manager ActionsManager
	bus     ActionsBus
}

func NewActionsModel(m ActionsManager, bus bus.Controller, log logging.Logger) *ActionsModel {
	return &ActionsModel{
		logger:  log,
		manager: m,
		bus:     ActionsBus{bus},
	}
}

func (m *ActionsModel) Add(a *Action) error {
	m.Actions.Store(a.ID, a)
	a.bus = m.bus
	a.updateStatus(ActionPending)
	switch a.Name {
	case plugin.ActionPair:
		timeout := to.Int(a.Input["timeout"])
		err := m.manager.AddNewThings(timeout)
		if err != nil {
			a.Error = err.Error()
			a.updateStatus(ActionError)
			m.logger.Infof("things was not added.err:%s", err.Error())
			return err
		}
		a.updateStatus(ActionCompleted)
	case plugin.ActionUnpair:

	}
	return nil
}

func (m *ActionsModel) OnActionStatus(model *Action) {

}
