package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type BooleanPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.BooleanSchema
	Observable bool `json:"observable,omitempty"`
}

func (p *BooleanPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.BooleanSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.BooleanSchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}

func (p *BooleanPropertyAffordance) MarshalJSON() ([]byte, error) {
	type property struct {
		*ia.InteractionAffordance
		*schema.BooleanSchema
		Observable bool `json:"observable,omitempty"`
	}
	prop := property{
		InteractionAffordance: p.InteractionAffordance,
		BooleanSchema:         p.BooleanSchema,
		Observable:            p.Observable,
	}
	return json.Marshal(prop)
}
