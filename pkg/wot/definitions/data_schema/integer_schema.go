package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type IntegerSchema struct {
	*DataSchema
	Minimum          controls.Integer `json:"minimum"`
	ExclusiveMinimum controls.Integer `json:"exclusiveMinimum,omitempty"`
	Maximum          controls.Integer `json:"maximum,omitempty"`
	ExclusiveMaximum controls.Integer `json:"exclusiveMaximum,omitempty"`
	MultipleOf       controls.Integer `json:"multipleOf,omitempty"`
}

func NewIntegerSchemaFromString(description string) *IntegerSchema {
	data := []byte(description)
	var schema = IntegerSchema{}
	schema.DataSchema = NewDataSchemaFromString(description)
	if schema.DataSchema == nil && schema.DataSchema.GetType() != controls.TypeString {
		return nil
	}
	schema.Minimum = controls.Integer(json.Get(data, "minimum").ToInt64())
	schema.ExclusiveMinimum = controls.Integer(json.Get(data, "exclusiveMinimum").ToInt64())
	schema.Maximum = controls.Integer(json.Get(data, "maximum").ToInt64())
	schema.ExclusiveMaximum = controls.Integer(json.Get(data, "exclusiveMaximum").ToInt64())
	schema.MultipleOf = controls.Integer(json.Get(data, "multipleOf").ToInt64())
	return &schema
}

func (i *IntegerSchema) Convert(v interface{}) interface{} {
	return i.clamp(controls.ToInteger(v))
}

func (i *IntegerSchema) clamp(value controls.Integer) controls.Integer {
	if i.Maximum != 0 {
		if value > i.Maximum {
			return i.Maximum
		}
	}
	if value < i.Minimum {
		return i.Minimum
	}
	return value
}

func (i IntegerSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(i)
}
