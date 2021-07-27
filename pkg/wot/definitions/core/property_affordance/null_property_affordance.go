package property_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)
import "github.com/galenliu/gateway/pkg/wot/definitions/core"

type NullPropertyAffordance struct {
	*core.InteractionAffordance
	*schema.NullSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

func NewNullPropertyAffordanceFormString(description string) *NullPropertyAffordance {
	data := []byte(description)
	var p = NullPropertyAffordance{}
	p.InteractionAffordance = core.NewInteractionAffordanceFromString(description)
	p.NullSchema = schema.NewNullSchemaFromString(description)
	if p.NullSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
