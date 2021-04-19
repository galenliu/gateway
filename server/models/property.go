package models

import (
	"addon/wot"
	json "github.com/json-iterator/go"
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

	var property Property
	property.SetAtType(gjson.Get(description, "@type").String())
	property.SetTitle(gjson.Get(description, "title").String())
	typ := gjson.Get(description, "type").String()
	if typ == "" {
		return nil
	}
	switch typ {

	case "array":
		var prop wot.ArraySchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	case "boolean ":
		var prop wot.BooleanSchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	case "number":
		var prop wot.NumberSchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	case "integer":
		var prop wot.IntegerSchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	case "object":
		var prop wot.ObjectSchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	case "string":
		var prop wot.StringSchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	case "null":
		var prop wot.NullSchema
		err := json.UnmarshalFromString(description, &prop)
		if err != nil {
			property.IDataSchema = prop
		}

	}
	return &property
}
