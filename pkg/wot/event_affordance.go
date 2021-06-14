package wot

import (
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type EventAffordance struct {
	*InteractionAffordance
	Subscription IDataSchema `json:"subscription,omitempty"`
	Data         IDataSchema `json:"data,omitempty"`
	Cancellation IDataSchema `json:"cancellation,omitempty"`
}

func NewEventAffordanceFromString(data string) *EventAffordance {
	var ia = InteractionAffordance{}
	err := json.Unmarshal([]byte(data), &ia)
	if err != nil {
		return nil
	}
	var e = EventAffordance{}

	if gjson.Get(data, "subscription").Exists() {
		s := gjson.Get(data, "subscription").String()
		d := NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	if gjson.Get(data, "data").Exists() {
		s := gjson.Get(data, "data").String()
		d := NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	if gjson.Get(data, "cancellation").Exists() {
		s := gjson.Get(data, "cancellation").String()
		d := NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	e.InteractionAffordance = &ia
	return &e
}
