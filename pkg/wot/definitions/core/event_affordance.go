package core

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type EventAffordance struct {
	ia.InteractionAffordance
	Subscription *schema.DataSchema `json:"subscription,omitempty"`
	Data         *schema.DataSchema `json:"data,omitempty"`
	Cancellation *schema.DataSchema `json:"cancellation,omitempty"`
}

func (e *EventAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	e.InteractionAffordance = i
	if err != nil {
		return err
	}
	var s schema.DataSchema
	json.Get(data, "subscription").ToVal(&s)
	e.Subscription = &s
	var d schema.DataSchema
	json.Get(data, "data").ToVal(&d)
	e.Subscription = &d
	var c schema.DataSchema
	json.Get(data, "cancellation").ToVal(&c)
	e.Subscription = &d
	return nil
}
