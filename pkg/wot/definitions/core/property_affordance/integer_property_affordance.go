package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type IntegerPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.IntegerSchema
	Observable bool             `json:"observable,omitempty"`
	Value      controls.Integer `json:"value,omitempty"`
}

func NewIntegerPropertyAffordanceFormString(description string) *IntegerPropertyAffordance {
	data := []byte(description)
	var p = IntegerPropertyAffordance{}
	p.InteractionAffordance = ia.NewInteractionAffordanceFromString(description)
	p.IntegerSchema = schema.NewIntegerSchemaFromString(description)
	if p.IntegerSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
