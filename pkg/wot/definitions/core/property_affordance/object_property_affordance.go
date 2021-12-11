package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type ObjectPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.ObjectSchema
	Observable bool `json:"observable,omitempty"`
}

func (p *ObjectPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.ObjectSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.ObjectSchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}

func (p *ObjectPropertyAffordance) MarshalJSON() ([]byte, error) {
	type property struct {
		*ia.InteractionAffordance
		*schema.ObjectSchema
		Observable bool `json:"observable,omitempty"`
	}
	prop := property{
		InteractionAffordance: p.InteractionAffordance,
		ObjectSchema:          p.ObjectSchema,
		Observable:            p.Observable,
	}
	return json.Marshal(prop)
}
