package property_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)
import "github.com/galenliu/gateway/pkg/wot/definitions/core"

type IntegerPropertyAffordance struct {
	*core.InteractionAffordance
	*schema.IntegerSchema
	Observable bool `json:"observable"`
}

func NewIntegerPropertyAffordanceFormString(description string) *IntegerPropertyAffordance {
	data := []byte(description)
	var p = IntegerPropertyAffordance{}
	p.InteractionAffordance = core.NewInteractionAffordanceFromString(description)
	p.IntegerSchema = schema.NewIntegerSchemaFromString(description)
	if p.IntegerSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
