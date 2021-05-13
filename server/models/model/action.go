package model

import (
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
	"strconv"
)

var actionId = 0

func generateActionId() string {
	actionId = actionId + 1
	return strconv.Itoa(actionId)
}

type Action struct {
	*wot.ActionAffordance
	ID      string `json:"id"`
	Href    string `json:"href,omitempty"`
	Name    string `json:"name,omitempty"`
	ThingId string
}

func NewActionFromString(data string) *Action {
	var this = Action{}
	aa := wot.NewActionAffordanceFromString(data)
	if aa.Forms == nil {
		aa.Forms = append(aa.Forms, wot.Form{
			Href: "",
			Op:   []string{wot.InvokeAction},
		})
	}

	this.ActionAffordance = aa
	return &this
}

func (action *Action) getId() string {
	return action.ID
}

func (action *Action) getInput() string {
	s, _ := json.MarshalToString(action.Input)
	return s
}
