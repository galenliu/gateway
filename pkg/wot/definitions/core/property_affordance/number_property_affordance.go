package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type NumberPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

func NewNumberPropertyAffordanceFormString(description string) *NumberPropertyAffordance {
	data := []byte(description)
	var p = NumberPropertyAffordance{}
	p.InteractionAffordance = ia.NewInteractionAffordanceFromString(description)
	p.NumberSchema = schema.NewNumberSchemaFromString(description)
	if p.NumberSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
