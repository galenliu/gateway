package models

import (
	"encoding/json"
	"github.com/galenliu/gateway-addon/wot"
	"github.com/tidwall/gjson"
)

type Property struct {
	*wot.PropertyAffordance
	Name string `json:"name"`

	Minimum interface{} `json:"minimum,omitempty"`
	Maximum interface{} `json:"maximum,omitempty"`
	ThingId string      `json:"thingId"`
}

func NewProperty(description string) *Property {

	var property = Property{PropertyAffordance: new(wot.PropertyAffordance)}

	var i wot.InteractionAffordance
	err := json.Unmarshal([]byte(description), &i)
	if err != nil {
		property.InteractionAffordance = &i
	}

	if !gjson.Get(description, "type").Exists() {
		return nil
	}
	typ := gjson.Get(description, "type").String()
	if typ == "" {
		return nil
	}
	switch typ {
	case "array":
		var prop wot.ArraySchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}

	case "boolean":
		var prop wot.BooleanSchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}

	case "number":
		var prop wot.NumberSchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}

	case "integer":
		var prop wot.IntegerSchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}

	case "object":
		var prop wot.ObjectSchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}

	case "string":
		var prop wot.StringSchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}

	case "null":
		var prop wot.NullSchema
		err := json.Unmarshal([]byte(description), &prop)
		if err == nil {
			property.IDataSchema = &prop
		}
	default:
		return nil
	}
	if !gjson.Get(description, "@type").Exists() {
		return nil
	}
	if !gjson.Get(description, "title").Exists() {
		property.SetTitle(gjson.Get(description, "title").String())
	}
	return &property
}
