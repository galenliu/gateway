package triggers

import (
	"github.com/galenliu/gateway/pkg/rules_engine/property"
)

type PropertyTriggerDescription struct {
	Property property.PropertyDescription `json:"property"`
}

type PropertyTrigger struct {
	*property.Property
}

func NewPropertyTrigger(des PropertyTriggerDescription, bus property.Bus, things property.ThingsHandler) *PropertyTrigger {
	return &PropertyTrigger{
		Property: property.NewProperty(des.Property, bus, things),
	}
}
