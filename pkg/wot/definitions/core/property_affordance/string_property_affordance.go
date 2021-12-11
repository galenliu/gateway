package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type StringPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.StringSchema
	Observable bool `json:"observable,omitempty"`
}

func (p *StringPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.StringSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.StringSchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}

func (p *StringPropertyAffordance) MarshalJSON() ([]byte, error) {
	type property struct {
		*ia.InteractionAffordance
		*schema.StringSchema
		Observable bool `json:"observable,omitempty"`
	}
	prop := property{
		InteractionAffordance: p.InteractionAffordance,
		StringSchema:          p.StringSchema,
		Observable:            p.Observable,
	}
	return json.Marshal(prop)
}
