package thing

import (
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
	"strconv"
	"time"
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
	ThingId       string
}

type Input struct {
	Timeout string `json:"timeout,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func NewAction(name string, input map[string]interface{}, th *Thing) *Action {
	a := &Action{
		ID:            generateActionId(),
		Name:          name,
		Input:         input,
		Status:        "created",
		TimeRequested: time.Now().String(),
	}
	if th != nil {
		a.ThingId = th.ID
		a.Href = fmt.Sprintf("%s/%s/%s/%s", th.ID, util.ActionsPath, name, a.ID)
	}
	return a
}

func NewThingAction(thingId, name string, input map[string]interface{}) *Action {
	a := NewAction(name, input, nil)
	a.ThingId = thingId
	a.Href = "/" + thingId + a.Href
	return a
}

func (action *Action) GetDescription() string {
	dsc, err := json.MarshalToString(action)
	if err != nil {
		return ""
	}
	return dsc
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
