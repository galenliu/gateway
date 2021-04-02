package thing

import json "github.com/json-iterator/go"

type Event struct {
	ID      string `json:"-" gorm:"primaryKey"`
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewEvent() *Event {
	e := &Event{}
	return e
}

func (e Event) GetThingId() string {
	return e.ThingId
}

func (e Event) GetDescription() string {
	data, _ := json.MarshalToString(e)
	return data

}
