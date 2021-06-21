package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type IntegerSchema struct {
	*dataSchema
	Minimum          controls.Integer `json:"minimum"`
	ExclusiveMinimum controls.Integer `json:"exclusiveMinimum,omitempty"`
	Maximum          controls.Integer `json:"maximum,omitempty"`
	ExclusiveMaximum controls.Integer `json:"exclusiveMaximum,omitempty"`
	MultipleOf       controls.Integer `json:"multipleOf,omitempty"`
}

func NewIntegerSchemaFromString(description string) *IntegerSchema {
	data := []byte(description)
	var schema = IntegerSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil && schema.dataSchema.GetType() != controls.TypeString {
		return nil
	}
	schema.Minimum = controls.Integer(json.Get(data, "minimum").ToInt64())
	schema.ExclusiveMinimum = controls.Integer(json.Get(data, "exclusiveMinimum").ToInt64())
	schema.Maximum = controls.Integer(json.Get(data, "maximum").ToInt64())
	schema.ExclusiveMaximum = controls.Integer(json.Get(data, "exclusiveMaximum").ToInt64())
	schema.MultipleOf = controls.Integer(json.Get(data, "multipleOf").ToInt64())
	return &schema
}

func (n IntegerSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
