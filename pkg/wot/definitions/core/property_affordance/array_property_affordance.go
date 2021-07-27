package property_affordance

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type ArrayPropertyAffordance struct {
	*core.InteractionAffordance
	*schema.ArraySchema
	Observable bool `json:"observable"`
}

func NewArrayPropertyAffordanceFormString(description string) *ArrayPropertyAffordance {
	data := []byte(description)
	var p = ArrayPropertyAffordance{}
	p.InteractionAffordance = core.NewInteractionAffordanceFromString(description)
	p.ArraySchema = schema.NewArraySchemaFromString(description)
	if p.ArraySchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
