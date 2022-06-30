package core

import (
	"encoding/json"
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	"github.com/tidwall/gjson"
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
	result := gjson.GetBytes(data, "subscription")
	err = json.Unmarshal(data[result.Index:result.Index+len(result.Raw)], &s)
	if err != nil {
		e.Subscription = &s
	}

	result = gjson.GetBytes(data, "data")
	err = json.Unmarshal(data[result.Index:result.Index+len(result.Raw)], &s)
	if err != nil {
		e.Data = &s
	}

	result = gjson.GetBytes(data, "data")
	err = json.Unmarshal(data[result.Index:result.Index+len(result.Raw)], &s)
	if err != nil {
		e.Data = &s
	}

	result = gjson.GetBytes(data, "cancellation")
	err = json.Unmarshal(data[result.Index:result.Index+len(result.Raw)], &s)
	if err != nil {
		e.Cancellation = &s
	}
	return nil
}
