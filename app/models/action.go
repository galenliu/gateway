package models

import (
	"gateway/event"
	json "github.com/json-iterator/go"
	"time"
)

const (
	Created   = "created"
	Completed = "completed"
	Error     = "error"
)

type Action struct {
	ID            int         `json:"-"`
	Name          string      `json:"name"`
	Input         interface{} `json:"input"`
	Href          string      `json:"href"`
	Status        string      `json:"status"`
	TimeRequested string      `json:"time_requested"`
	TimeCompleted string      `json:"time_completed,omitempty"`
	Error         string      `json:"error,omitempty"`

	ThingID string `json:"-"`
}

type Input struct {
	Timeout string `json:"timeout,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func NewAction(name string, input interface{}) *Action {
	nextId := generateId()
	return &Action{
		ID:            nextId,
		Name:          name,
		Input:         input,
		Href:          "", //fmt.Sprintf("%s/%s/%v", app.ActionsPath, name, nextId),
		Status:        "created",
		TimeRequested: time.Now().String(),
	}
}

func NewThingAction(thingId, actionName string, input interface{}) *Action {
	return &Action{
		ID:            0,
		Name:          "",
		Input:         nil,
		Href:          "",
		Status:        "",
		ThingID:       thingId,
		TimeRequested: time.Now().String(),
		TimeCompleted: "",
		Error:         "",
	}

}

func (action *Action) GetDescription() ([]byte, error) {
	var desc = make(map[string]*Action)
	desc[action.Name] = action
	return json.Marshal(desc)
}

func (action *Action) updateStatus(newStatus string) {
	if action.Status == newStatus {
		return
	}
	if newStatus == "completed" {
		action.TimeCompleted = time.Now().String()
	}
    event.FireAction(action)
	action.Status = newStatus
}
