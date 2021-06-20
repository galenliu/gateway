package wot

import (
	"github.com/galenliu/gateway/wot/definitions/core"
	"github.com/tidwall/gjson"
)

type Property struct {
	core.PropertyAffordance
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewPropertyFromString(description string) *Property {
	var property = Property{}
	property.PropertyAffordance = core.NewPropertyAffordanceFromString(description)
	if gjson.Get(description, "name").Exists() {
		property.Name = gjson.Get(description, "name").String()
	}
	if gjson.Get(description, "thingId").Exists() {
		property.ThingId = gjson.Get(description, "thingId").String()
	}
	return &property
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetThingId() string {
	return p.ThingId
}
