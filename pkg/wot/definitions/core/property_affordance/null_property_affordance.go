package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type NullPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NullSchema
	Observable bool `json:"observable"`
}

func (p *NullPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.NullSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.NullSchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}

func (p *NullPropertyAffordance) MarshalJSON() ([]byte, error) {
	type property struct {
		*ia.InteractionAffordance
		*schema.NullSchema
		Observable bool `json:"observable,omitempty"`
	}
	prop := property{
		InteractionAffordance: p.InteractionAffordance,
		NullSchema:            p.NullSchema,
		Observable:            p.Observable,
	}
	return json.Marshal(prop)
}
