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

	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	schema.DataSchema = &dataSchema

	if schema.DataSchema == nil && schema.DataSchema.GetType() != controls.TypeInteger {
		return fmt.Errorf("type must integer")
	}
	schema.Minimum = controls.JSONGetInteger(data, "minimum")
	schema.ExclusiveMinimum = controls.JSONGetInteger(data, "exclusiveMinimum")
	schema.Maximum = controls.JSONGetInteger(data, "maximum")
	schema.ExclusiveMaximum = controls.JSONGetInteger(data, "exclusiveMaximum")
	schema.MultipleOf = controls.JSONGetInteger(data, "multipleOf")
	return nil
}

func (schema *IntegerSchema) GetDefaultValue() interface{} {
	if schema.DataSchema.Default != nil {
		return schema.Convert(schema.Default)
	}
	if len(schema.Enum) > 0 {
		return schema.Convert(schema.Enum[0])
	}
	if schema.Minimum != nil {
		return schema.Convert(*schema.Minimum)
	}
	return controls.DefaultInteger
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

func (schema *IntegerSchema) verify(value interface{}) bool {
	switch value.(type) {
	case controls.Integer:
		return true
	}
	return false
}
