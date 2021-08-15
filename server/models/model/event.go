package model

import (
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type Event struct {
	wot.EventAffordance

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
			ContentType: schema.ApplicationJson,
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
