package model

import (
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"strconv"
)

var actionId = 0

func generateActionId() string {
	actionId = actionId + 1
	return strconv.Itoa(actionId)
}

type Action struct {
	ID      string `json:"id"`
	Href    string `json:"href,omitempty"`
	Name    string `json:"name,omitempty"`
	ThingId string
}

func NewActionFromString(data string) *Action {
	var this = Action{}
	aa := wot.NewActionAffordanceFromString(data)
	if aa.Forms == nil {
		aa.Forms = append(aa.Forms, hypermedia_controls.Form{
			Href: "",
		})
	}

	return &this
}

func (action *Action) getId() string {
	return action.ID
}

func (action *Action) getInput() string {

	return ""
}
