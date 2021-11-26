package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
	uuid "github.com/satori/go.uuid"
	"time"
)

type bus interface {
	PublishActionStatus(action interface{})
}

type ActionDescription struct {
	Input         map[string]interface{} `json:"input,omitempty"`
	Href          string                 `json:"href,omitempty"`
	TimeRequested *time.Time             `json:"timeRequested,omitempty"`
	TimeCompleted *time.Time             `json:"timeCompleted,omitempty"`
	Error         error                  `json:"error,omitempty"`
}

type Action struct {
	Id            string                 `json:"-"`
	ThingId       string                 `json:"-"`
	Name          string                 `json:"name"`
	Input         map[string]interface{} `json:"input,omitempty"`
	Href          string                 `json:"href"`
	TimeRequested *time.Time             `json:"timeRequested,omitempty"`
	TimeCompleted *time.Time             `json:"timeCompleted,omitempty"`
	Status        string                 `json:"status"`
	Error         error                  `json:"error,omitempty"`
	bus           bus
	logger        logging.Logger
}

func NewActionModel(name string, input map[string]interface{}, thing *container.Thing, bus bus, log logging.Logger) *Action {
	t := time.Now()
	a := &Action{
		logger:        log,
		Input:         input,
		Id:            uuid.NewV4().String(),
		Name:          name,
		Status:        ActionCreated,
		TimeRequested: &t,
		TimeCompleted: nil,
		bus:           bus,
	}
	if thing != nil {
		a.ThingId = thing.GetId()
		a.Href = fmt.Sprintf("%s/%s/%s/%s", thing.GetHref(), constant.ActionsPath, name, a.Id)
	} else {
		a.Href = fmt.Sprintf("%s/%s/%s", constant.ActionsPath, name, a.Id)
	}
	return a
}

func (action *Action) updateStatus(newStatus string) {
	if action.Status == newStatus {
		return
	}
	if newStatus == ActionCompleted {
		t := time.Now()
		action.TimeCompleted = &t
	}
	action.Status = newStatus
	action.bus.PublishActionStatus(action)
}

func (action *Action) GetDescription() *ActionDescription {
	return &ActionDescription{
		Input:         action.Input,
		Href:          action.Href,
		TimeRequested: action.TimeRequested,
		TimeCompleted: action.TimeRequested,
		Error:         nil,
	}
}

func (action *Action) GetId() string {
	return action.Id
}

func (action *Action) GetName() string {
	return action.Name
}

func (action *Action) GetInput() map[string]interface{} {
	return action.Input
}

func (action *Action) SetErr(err error) {
	action.Error = err
	action.updateStatus(ActionError)
}
