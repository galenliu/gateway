package property_affordance

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type StringPropertyAffordance struct {
	*interaction_affordance.InteractionAffordance
	*schema.StringSchema
	Observable bool `json:"observable"`
}

func NewStringPropertyAffordanceFormString(description string) *StringPropertyAffordance {
	data := []byte(description)
	var p = StringPropertyAffordance{}
	p.InteractionAffordance = interaction_affordance.NewInteractionAffordanceFromString(description)
	p.StringSchema = schema.NewStringSchemaFromString(description)
	if p.StringSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
