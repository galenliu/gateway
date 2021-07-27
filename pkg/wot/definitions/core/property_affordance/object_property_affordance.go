package property_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)
import "github.com/galenliu/gateway/pkg/wot/definitions/core"

type ObjectPropertyAffordance struct {
	*core.InteractionAffordance
	*schema.ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

func NewObjectPropertyAffordanceFormString(description string) *ObjectPropertyAffordance {
	data := []byte(description)
	var p = ObjectPropertyAffordance{}
	p.InteractionAffordance = core.NewInteractionAffordanceFromString(description)
	p.ObjectSchema = schema.NewObjectSchemaFromString(description)
	if p.ObjectSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
