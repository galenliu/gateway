package core

import (
	data_schema2 "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	hypermedia_controls2 "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/galenliu/gateway/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type EventAffordance interface {
}

type eventAffordance struct {
	*InteractionAffordance
	Subscription data_schema.dataSchema `json:"subscription,omitempty"`
	Data         data_schema.dataSchema `json:"data,omitempty"`
	Cancellation data_schema.dataSchema `json:"cancellation,omitempty"`
}

func NewEventAffordanceFromString(data string) EventAffordance {
	var ia = InteractionAffordance{}
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
		e.Forms = append(e.Forms, hypermedia_controls2.Form{
			Href:        "",
			ContentType: data_schema2.ApplicationJson,
			Op:          []string{hypermedia_controls2.SubscribeEvent},
		})
	}
	e.InteractionAffordance = &ia
	return e
}
