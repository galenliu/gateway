package wot

import (
	"github.com/galenliu/gateway/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type Event struct {
	core.EventAffordance
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewEventFromString(data string) *Event {
	var this = Event{}
	aa := core.NewEventAffordanceFromString(data)
	this.EventAffordance = aa
	return &this
}

func (e Event) GetThingId() string {
	return e.ThingId
}

func (e Event) ToJson() string {
	data, _ := json.MarshalToString(e)
	return data

}
