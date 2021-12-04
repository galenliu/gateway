package data_schema

import (
	"fmt"
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

func (schema *IntegerSchema) UnmarshalJSON(data []byte) error {

	getInteger := func(sep string) *controls.Integer {
		var ig controls.Integer
		json.Get(data, sep).ToVal(&ig)
		if &ig != nil {
			return &ig
		}
		return nil
	}
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	schema.DataSchema = &dataSchema

	if schema.DataSchema == nil && schema.DataSchema.GetType() != controls.TypeInteger {
		return fmt.Errorf("type must integer")
	}
	schema.Minimum = getInteger("minimum")
	schema.ExclusiveMinimum = getInteger("exclusiveMinimum")
	schema.Maximum = getInteger("maximum")
	schema.ExclusiveMaximum = getInteger("exclusiveMaximum")
	schema.MultipleOf = getInteger("multipleOf")
	return nil
}

func (schema *IntegerSchema) GetDefaultValue() interface{} {
	if schema.DataSchema.Default != nil {
		return schema.Convert(schema.Default)
	}
	return nil
}

func (schema *IntegerSchema) Convert(v interface{}) interface{} {
	return schema.clamp(controls.ToInteger(v))
}

func (schema *IntegerSchema) clamp(value controls.Integer) controls.Integer {
	if schema.Maximum != nil {
		if value > *schema.Maximum {
			return *schema.Maximum
		}
	}
	if schema.Minimum != nil {
		if value < *schema.Minimum {
			return *schema.Minimum
		}

	}
	return value
}
