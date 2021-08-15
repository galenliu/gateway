package model

import (
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
)

type Property struct {
	wot.PropertyAffordance
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewPropertyFromString(description string) *Property {
	//data := []byte(description)
	return nil
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetThingId() string {
	return p.ThingId
}
