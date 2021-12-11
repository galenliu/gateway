package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type IntegerPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.IntegerSchema
	Observable bool `json:"observable,omitempty"` //with default

}

func (p *IntegerPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.IntegerSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.IntegerSchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}

func (p *IntegerPropertyAffordance) MarshalJSON() ([]byte, error) {
	type property struct {
		*ia.InteractionAffordance
		*schema.IntegerSchema
		Observable bool `json:"observable,omitempty"`
	}
	prop := property{
		InteractionAffordance: p.InteractionAffordance,
		IntegerSchema:         p.IntegerSchema,
		Observable:            p.Observable,
	}
	return json.Marshal(prop)
}
