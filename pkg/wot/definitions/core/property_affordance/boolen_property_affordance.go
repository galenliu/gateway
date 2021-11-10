package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type BooleanPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

func (p BooleanPropertyAffordance) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&p, "", "   ")
}

func (p BooleanPropertyAffordance) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p)
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
