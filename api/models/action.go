package models

import (
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ActionDescription struct {
	Input         map[string]any `json:"input,omitempty"`
	Href          string         `json:"href,omitempty"`
	TimeRequested *time.Time     `json:"timeRequested,omitempty"`
	TimeCompleted *time.Time     `json:"timeCompleted,omitempty"`
	Error         error          `json:"error,omitempty"`
}

type Action struct {
	Id            string         `json:"-"`
	ThingId       string         `json:"-"`
	Name          string         `json:"name,omitempty"`
	Input         map[string]any `json:"input,omitempty"`
	Href          string         `json:"href,omitempty"`
	TimeRequested *time.Time     `json:"timeRequested,omitempty"`
	TimeCompleted *time.Time     `json:"timeCompleted,omitempty"`
	Status        string         `json:"status,omitempty"`
	Error         error          `json:"error,omitempty"`
	logger        logging.Logger
}

func NewActionModel(name string, input map[string]any, log logging.Logger, things ...*container.Thing) *Action {
	t := time.Now()
	a := &Action{
		logger:        log,
		Input:         input,
		Id:            uuid.NewV4().String(),
		Name:          name,
		TimeRequested: &t,
		TimeCompleted: nil,
		ThingId:       "",
	}
	a.updateStatus(ActionCreated)
	if things != nil && len(things) > 0 {
		thing := things[0]
		a.ThingId = thing.GetId()
		a.Href = fmt.Sprintf("%s/%s/%s/%s", thing.GetHref(), constant.ActionsPath, name, a.Id)
	} else {
		a.Href = fmt.Sprintf("%s/%s/%s", constant.ActionsPath, name, a.GetId())
	}
	return a
}

func (action *Action) GetDescription() ActionDescription {
	des := ActionDescription{
		Input:         action.Input,
		Href:          action.Href,
		TimeRequested: action.TimeRequested,
		TimeCompleted: action.TimeRequested,
		Error:         action.Error,
	}
	return des
}

func (action *Action) GetId() string {
	return action.Id
}

func (action *Action) GetThingId() string {
	return action.ThingId
}

func (action *Action) GetName() string {
	return action.Name
}

func (action *Action) GetStatus() string {
	return action.Status
}

func (action *Action) GetInput() map[string]any {
	return action.Input
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
	action.logger.Infof("actions.updateStatus: %s", newStatus)
}

func (action *Action) SetErr(err error) {
	action.Error = err
	action.updateStatus(ActionError)
}

func (action *Action) update(ad *actions.ActionDescription) {
	t, _ := time.Parse("2006-1-2 15:04:05", ad.TimeRequested)
	action.TimeRequested = &t
	if ad.TimeCompleted != "" {
		t, _ := time.Parse("2006-1-2 15:04:05", ad.TimeCompleted)
		action.TimeCompleted = &t
	}
	action.updateStatus(ad.Status)
}
