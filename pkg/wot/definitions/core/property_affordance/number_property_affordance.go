package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type NumberPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

func (p NumberPropertyAffordance) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&p, "", "   ")
}

func (p NumberPropertyAffordance) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p)
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
