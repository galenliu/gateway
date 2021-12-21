package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type IntegerSchema struct {
	*DataSchema
	Minimum          *controls.Integer `json:"minimum,,omitempty"`
	ExclusiveMinimum *controls.Integer `json:"exclusiveMinimum,omitempty"`
	Maximum          *controls.Integer `json:"maximum,omitempty"`
	ExclusiveMaximum *controls.Integer `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *controls.Integer `json:"multipleOf,omitempty"`
}

func (schema *IntegerSchema) GetDefaultValue() any {
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

func (schema *IntegerSchema) Convert(v any) any {
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

func (schema *IntegerSchema) verify(value any) bool {
	switch value.(type) {
	case controls.Integer:
		return true
	}
	return false
}
