package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type ArrayPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.ArraySchema
	Observable bool `json:"observable,omitempty"`
}

func (p *ArrayPropertyAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &i
	var s schema.ArraySchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.ArraySchema = &s
	p.Observable = json.Get(data, "observable").ToBool()
	return nil
}
