package core

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type EventAffordance interface {
}

type eventAffordance struct {
	*interaction_affordance.InteractionAffordance
	Subscription *data_schema.DataSchema `json:"subscription,omitempty"`
	Data         *data_schema.DataSchema `json:"data,omitempty"`
	Cancellation *data_schema.DataSchema `json:"cancellation,omitempty"`
}

func NewEventAffordanceFromString(data string) *eventAffordance {
	var ia = interaction_affordance.InteractionAffordance{}
	err := json.Unmarshal([]byte(data), &ia)
	if err != nil {
		return nil
	}
	var e = eventAffordance{}

	if gjson.Get(data, "subscription").Exists() {
		s := gjson.Get(data, "subscription").String()
		d := data_schema.NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	if gjson.Get(data, "data").Exists() {
		s := gjson.Get(data, "data").String()
		d := data_schema.NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	if gjson.Get(data, "cancellation").Exists() {
		s := gjson.Get(data, "cancellation").String()
		d := data_schema.NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}
	if e.Forms == nil {
		e.Forms = append(e.Forms, controls.Form{
			Href:        "",
			ContentType: data_schema.ApplicationJson,
			Op:          []string{controls.SubscribeEvent},
		})
	}
	e.InteractionAffordance = &ia
	return &e
}
