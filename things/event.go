package things

import (
	core2 "github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type Event struct {
	core2.EventAffordance
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewEventFromString(data string) *Event {
	var this = Event{}
	aa := core2.NewEventAffordanceFromString(data)
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
