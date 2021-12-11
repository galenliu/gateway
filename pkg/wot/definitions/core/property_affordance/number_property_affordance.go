package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type NumberPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NumberSchema
	Observable bool `json:"observable,omitempty"`
}

func (p *NumberPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.NumberSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.NumberSchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}

func (p *NumberPropertyAffordance) MarshalJSON() ([]byte, error) {
	type property struct {
		*ia.InteractionAffordance
		*schema.NumberSchema
		Observable bool `json:"observable,omitempty"`
	}
	prop := property{
		InteractionAffordance: p.InteractionAffordance,
		NumberSchema:          p.NumberSchema,
		Observable:            p.Observable,
	}
	return json.Marshal(prop)
}
