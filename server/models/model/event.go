package model

import (
	"github.com/galenliu/gateway-addon/wot"
	"github.com/galenliu/gateway/wot/definitions/data_schema"
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
		aa.Forms = append(aa.Forms, hypermedia_controls.Form{
			Href:        "",
			ContentType: data_schema.ApplicationJson,
			Op:          []string{hypermedia_controls.SubscribeEvent},
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
