package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type IntegerPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.IntegerSchema
	Observable bool `json:"observable"`
}

func (p IntegerPropertyAffordance) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&p, "", "   ")
}

func (p IntegerPropertyAffordance) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p)
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
