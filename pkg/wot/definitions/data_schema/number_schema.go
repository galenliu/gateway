package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type NumberSchema struct {
	*DataSchema
	Minimum          *controls.Double `json:"minimum,omitempty"`
	ExclusiveMinimum *controls.Double `json:"exclusiveMinimum,omitempty"`
	Maximum          *controls.Double `json:"maximum,omitempty"`
	ExclusiveMaximum *controls.Double `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *controls.Double `json:"multipleOf,omitempty"`
}

func (schema *NumberSchema) GetDefaultValue() any {
	if schema.DataSchema.Default != nil {
		return schema.Default
	}
	if len(schema.Enum) > 0 {
		return schema.Enum[0]
	}
	if schema.Minimum != nil {
		return *schema.Minimum
	}
	return controls.DefaultNumber
}

func (schema *NumberSchema) Convert(v any) any {
	//return schema.clamp(controls.Double(to.Float64(v)))
	return v
}

func (schema NumberSchema) clamp(value controls.Double) controls.Double {
	if schema.Maximum != nil {
		if value > *schema.Maximum {
			return *schema.Maximum
		}
	}
	if schema.Maximum != nil {
		if value < *schema.Minimum {
			return *schema.Minimum
		}
	}
	return value
}

func (schema *NumberSchema) verify(value any) bool {
	switch value.(type) {
	case controls.Number:
		return true
	}
	return false
}
