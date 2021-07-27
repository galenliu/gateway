package property_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)
import "github.com/galenliu/gateway/pkg/wot/definitions/core"

type BooleanPropertyAffordance struct {
	*core.InteractionAffordance
	*schema.BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

func NewBooleanPropertyAffordanceFormString(description string) *BooleanPropertyAffordance {
	data := []byte(description)
	var p = BooleanPropertyAffordance{}
	p.InteractionAffordance = core.NewInteractionAffordanceFromString(description)
	p.BooleanSchema = schema.NewBooleanSchemaFromString(description)
	if p.BooleanSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
