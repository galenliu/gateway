package data_schema

import (
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

func NewNumberSchemaFromString(description string) *NumberSchema {
	data := []byte(description)
	var schema = NumberSchema{}
	schema.DataSchema = NewDataSchemaFromString(description)
	if schema.DataSchema == nil || schema.DataSchema.GetType() != controls.TypeNumber {
		return nil
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
	return &schema
}

func (n *NumberSchema) Convert(v interface{}) interface{} {
	return n.clamp(controls.Double(to.Float64(v)))
}

func (n NumberSchema) clamp(value controls.Double) controls.Double {
	if n.Maximum != nil {
		if value > *n.Maximum {
			return *n.Maximum
		}
	}
	if n.Maximum != nil {
		if value < *n.Minimum {
			return *n.Minimum
		}
	}
	return value
}
