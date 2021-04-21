package models

import (
	"github.com/galenliu/gateway-addon/wot"
	json "github.com/json-iterator/go"
)

type Event struct {
	*wot.EventAffordance

	ID      string `json:"-" gorm:"primaryKey"`
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewEventFromString(data string) *Event {
	var this = Event{}
	aa := wot.NewEventAffordanceFromString(data)
	if aa.Forms == nil {
		aa.Forms = append(aa.Forms, wot.Form{
			Href:        "",
			ContentType: wot.ApplicationJson,
			Op:          []string{wot.SubscribeEvent},
		})
	}
	this.EventAffordance = aa
	return &this
}

func (e Event) GetThingId() string {
	return e.ThingId
}

func (e Event) GetDescription() string {
	data, _ := json.MarshalToString(e)
	return data

}
