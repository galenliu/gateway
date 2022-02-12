package triggers

import (
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/property"
)

type PropertyTriggerDescription struct {
	TriggerDescription
	Property property.Description `json:"property"`
}

type PropertyTrigger struct {
	*Trigger
	property *property.Property
}

func NewPropertyTrigger(des PropertyTriggerDescription, container container.Container) *PropertyTrigger {
	p := &PropertyTrigger{
		Trigger:  NewTrigger(des.TriggerDescription),
		property: property.NewProperty(des.Property, container),
	}
	return p
}

func (p *PropertyTrigger) Start() {
	p.property.Start()
}

func (p *PropertyTrigger) Stop() {
	p.property.Stop()
}

func (p *PropertyTrigger) OnValueChanged(value any) {
	fmt.Println("on value changed function not implemented")
}

func (p *PropertyTrigger) ToDescription() PropertyTriggerDescription {
	desc := PropertyTriggerDescription{}
	desc.TriggerDescription = p.Trigger.ToDescription()
	desc.Property = p.property.ToDescription()
	return desc
}

func (p *PropertyTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToDescription())
}
