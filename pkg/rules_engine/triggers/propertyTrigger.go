package triggers

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/property"
)

type PropertyTriggerDescription struct {
	Property property.PropertyDescription `json:"property"`
}

type PropertyTrigger struct {
	*Trigger
	property *property.Property
}

func NewPropertyTrigger(des PropertyTriggerDescription, container container.Container) *PropertyTrigger {
	p := &PropertyTrigger{
		property: property.NewProperty(des.Property, container),
	}
	return p
}
