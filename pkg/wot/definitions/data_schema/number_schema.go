package data_schema

import (
	"fmt"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
)

type NumberSchema struct {
	*DataSchema
	Minimum          *controls.Double `json:"minimum,omitempty"`
	ExclusiveMinimum *controls.Double `json:"exclusiveMinimum,omitempty"`
	Maximum          *controls.Double `json:"maximum,omitempty"`
	ExclusiveMaximum *controls.Double `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *controls.Double `json:"multipleOf,omitempty"`
}

func (schema *NumberSchema) UnmarshalJSON(data []byte) error {
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	if schema.DataSchema == nil || schema.DataSchema.GetType() != controls.TypeNumber {
		return fmt.Errorf("type must number")
	}
	getDouble := func(sep string) *controls.Double {
		var d controls.Double
		json.Get(data, sep).ToVal(&d)
		if &d != nil {
			return &d
		}
		return nil
	}
	schema.Minimum = getDouble("minimum")
	schema.Maximum = getDouble("maximum")
	schema.ExclusiveMinimum = getDouble("exclusiveMinimum")
	schema.ExclusiveMaximum = getDouble("exclusiveMaximum")
	schema.MultipleOf = getDouble("multipleOf")
	return nil
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
	return schema.clamp(controls.Double(to.Float64(v)))
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
