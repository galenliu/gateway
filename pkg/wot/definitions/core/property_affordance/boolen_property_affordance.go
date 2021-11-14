package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type BooleanPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.BooleanSchema
	Observable bool `json:"observable,omitempty"`
	Value      bool `json:"value,omitempty"`
}

func NewBooleanPropertyAffordanceFormString(description string) *BooleanPropertyAffordance {
	data := []byte(description)
	var p = BooleanPropertyAffordance{}
	p.InteractionAffordance = ia.NewInteractionAffordanceFromString(description)
	p.BooleanSchema = schema.NewBooleanSchemaFromString(description)
	if p.BooleanSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
