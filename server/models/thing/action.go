package thing

import (
	json "github.com/json-iterator/go"
	"time"
)

const (
	Created   = "created"
	Completed = "completed"
	Error     = "error"
)

type Action struct {
	ID            uint                   `json:"-"`
	Name          string                 `json:"name"`
	Input         map[string]interface{} `json:"input"`
	Href          string                 `json:"href"`
	Status        string                 `json:"status"`
	TimeRequested string                 `json:"time_requested"`
	TimeCompleted string                 `json:"time_completed,omitempty"`
	Error         string                 `json:"error,omitempty"`

	ActionID string `json:"-" gorm:"primaryKey"`
	ThingID  string
}

type Input struct {
	Timeout string `json:"timeout,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func NewAction(name string, input map[string]interface{}) *Action {

	return &Action{
		ID:            1,
		Name:          name,
		Input:         input,
		Href:          "",
		Status:        "created",
		TimeRequested: time.Now().String(),
	}
}

func NewThingAction(thingId, actionName string, input map[string]interface{}) *Action {
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
	action.Status = newStatus
}
