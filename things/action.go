package things

import (
	core2 "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"strconv"
)

var actionId = 0

func generateActionId() string {
	actionId = actionId + 1
	return strconv.Itoa(actionId)
}

type Action struct {
	core2.ActionAffordance
	ID      string `json:"id"`
	Href    string `json:"href,omitempty"`
	Name    string `json:"name,omitempty"`
	ThingId string
}

func NewActionFromString(data string) *Action {
	var action = Action{}
	a := core2.NewActionAffordanceFromString(data)
	action.ActionAffordance = a
	return &action
}

func (action *Action) getId() string {
	return action.ID
}
