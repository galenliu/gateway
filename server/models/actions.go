package models

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin"
	"github.com/xiam/to"
	"sync"
)

type ActionsManager interface {
	AddNewThings(timeOut int) error
}

type ActionsModel struct {
	logger  logging.Logger
	Actions sync.Map
	manager ActionsManager
}

func NewActionsModel(m ActionsManager, log logging.Logger) *ActionsModel {
	return &ActionsModel{
		logger:  log,
		manager: m,
	}
}

func (m *ActionsModel) Add(a *Action) error {
	m.Actions.Store(a.ID, a)
	switch a.Name {
	case plugin.ActionPair:
		timeout := to.Int(a.Input["timeout"])
		err := m.manager.AddNewThings(timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ActionsModel) OnActionStatus(model *Action) {

}
