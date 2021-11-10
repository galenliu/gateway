package core

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type EventAffordance interface {
}

type eventAffordance struct {
	*ia.InteractionAffordance
	Subscription *schema.DataSchema `json:"subscription,omitempty"`
	Data         *schema.DataSchema `json:"data,omitempty"`
	Cancellation *schema.DataSchema `json:"cancellation,omitempty"`
}

func NewEventAffordanceFromString(data string) *eventAffordance {
	var affordance = ia.InteractionAffordance{}
	err := json.Unmarshal([]byte(data), &affordance)
	if err != nil {
		return nil
	}
	var e = eventAffordance{}

	if gjson.Get(data, "subscription").Exists() {
		s := gjson.Get(data, "subscription").String()
		d := schema.NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	if gjson.Get(data, "data").Exists() {
		s := gjson.Get(data, "data").String()
		d := schema.NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}

	if gjson.Get(data, "cancellation").Exists() {
		s := gjson.Get(data, "cancellation").String()
		d := schema.NewDataSchemaFromString(s)
		if d != nil {
			e.Subscription = d
		}
	}
	if e.Forms == nil {
		e.Forms = append(e.Forms, controls.Form{
			Href:        "",
			ContentType: schema.ApplicationJson,
			Op:          controls.NewArrayOfString(controls.SubscribeEvent),
		})
	}
	e.InteractionAffordance = &affordance
	return &e
}
