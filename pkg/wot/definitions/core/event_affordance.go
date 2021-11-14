package core

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/tidwall/gjson"

	json "github.com/json-iterator/go"
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
	var affordance = ia.NewInteractionAffordanceFromString(data)
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

	s := json.Get([]byte(data), "data").ToString()
	d := schema.NewDataSchemaFromString(s)
	if d != nil {
		e.Subscription = d
	}

	s = json.Get([]byte(data), "cancellation").ToString()
	d = schema.NewDataSchemaFromString(s)
	if d != nil {
		e.Subscription = d
	}

	if e.Forms == nil {
		e.Forms = append(e.Forms, controls.Form{
			Href:        "",
			ContentType: schema.ApplicationJson,
			Op:          controls.NewArrayOfString(controls.SubscribeEvent),
		})
	}
	e.InteractionAffordance = affordance
	return &e
}
