package thing

import (
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
	"strconv"
	"time"
)

const (
	Created   = "created"
	Completed = "completed"
	Error     = "error"
)

var actionId int = 0

func generateActionId() string {
	actionId = actionId + 1
	return strconv.Itoa(actionId)
}

type Action struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Input         map[string]interface{} `json:"input"`
	Href          string                 `json:"href"`
	Status        string                 `json:"status"`
	TimeRequested string                 `json:"time_requested"`
	TimeCompleted string                 `json:"time_completed,omitempty"`
	Error         string                 `json:"error,omitempty"`
	ThingID       string
}

type Input struct {
	Timeout string `json:"timeout,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func NewAction(name string, input map[string]interface{}) *Action {
	a := &Action{
		ID:            generateActionId(),
		Name:          name,
		Input:         input,
		Status:        "created",
		TimeRequested: time.Now().String(),
	}
	a.Href = fmt.Sprintf("%s/%s/%s", util.ActionsPath, name, a.ID)
	return a
}

func NewThingAction(thingId, name string, input map[string]interface{}) *Action {
	a := NewAction(name, input)
	a.ThingID = thingId
	a.Href = "/" + thingId + a.Href
	return a
}

func (action *Action) GetDescription() (string, error) {
	dsc, err := json.MarshalToString(action)
	if err != nil {
		return "", err
	}
	return dsc, nil
}

func (action *Action) UpdateStatus(newStatus string) {
	if action.Status == newStatus {
		return
	}
	if newStatus == "completed" {
		action.TimeCompleted = time.Now().String()
	}
	action.Status = newStatus
	bus.Publish(util.ActionStatus, action)
}

func (action *Action) Update(a *Action) {
	action.TimeCompleted = a.TimeCompleted
	action.TimeRequested = a.TimeRequested
	if action != a {
		bus.Publish(util.ActionStatus, action)
	}

	bus.Publish(util.ActionStatus, action)
}

func (action *Action) getId() string {
	return action.ID
}

func (action *Action) getInput() string {
	s, _ := json.MarshalToString(action.Input)
	return s
}

func (action *Action) getTimeRequested() string {
	return action.TimeRequested
}

func (action *Action) getTimeCompleted() string {
	return action.TimeCompleted
}

func (action *Action) setErr(err error) {
	action.Status = err.Error()
	bus.Publish(util.ActionStatus, action)
}
