package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type IntegerSchema struct {
	*DataSchema
	Minimum          *controls.Integer `json:"minimum,,omitempty"`
	ExclusiveMinimum *controls.Integer `json:"exclusiveMinimum,omitempty"`
	Maximum          *controls.Integer `json:"maximum,omitempty"`
	ExclusiveMaximum *controls.Integer `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *controls.Integer `json:"multipleOf,omitempty"`
}

func NewIntegerSchemaFromString(description string) *IntegerSchema {
	data := []byte(description)
	var schema = IntegerSchema{}
	getInteger := func(sep string) *controls.Integer {
		var ig controls.Integer
		json.Get(data, sep).ToVal(&ig)
		if &ig != nil {
			return &ig
		}
		return nil
	}
	schema.DataSchema = NewDataSchemaFromString(description)
	if schema.DataSchema == nil && schema.DataSchema.GetType() != controls.TypeString {
		return nil
	}
	schema.Minimum = getInteger("minimum")
	schema.ExclusiveMinimum = getInteger("exclusiveMinimum")
	schema.Maximum = getInteger("maximum")
	schema.ExclusiveMaximum = getInteger("exclusiveMaximum")
	schema.MultipleOf = getInteger("multipleOf")
	return &schema
}

func (i *IntegerSchema) Convert(v interface{}) interface{} {
	return i.clamp(controls.ToInteger(v))
}

func (i *IntegerSchema) clamp(value controls.Integer) controls.Integer {
	if i.Maximum != nil {
		if value > *i.Maximum {
			return *i.Maximum
		}
	}
	if i.Minimum != nil {
		if value < *i.Minimum {
			return *i.Minimum
		}

	}
	return value
}
