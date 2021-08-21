package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
)

type NumberSchema struct {
	*dataSchema
	Minimum          controls.Double `json:"minimum"`
	ExclusiveMinimum controls.Double `json:"exclusiveMinimum,omitempty"`
	Maximum          controls.Double `json:"maximum,omitempty"`
	ExclusiveMaximum controls.Double `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64         `json:"multipleOf,omitempty"`
}

func NewNumberSchemaFromString(description string) *NumberSchema {
	data := []byte(description)
	var schema = NumberSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil || schema.dataSchema.GetType() != controls.TypeNumber {
		return nil
	}

	schema.Minimum = controls.Double(controls.ToNumber(json.Get(data, "minimum").ToFloat64()))
	schema.ExclusiveMinimum = controls.Double(json.Get(data, "exclusiveMinimum").ToFloat64())
	schema.Maximum = controls.Double(controls.ToNumber(json.Get(data, "maximum").ToFloat64()))
	schema.ExclusiveMaximum = controls.Double(controls.ToNumber(json.Get(data, "exclusiveMaximum").ToFloat64()))
	schema.MultipleOf = json.Get(data, "multipleOf").ToFloat64()
	return &schema
}

func (n *NumberSchema) Convert(v interface{}) interface{} {
	return n.clamp(controls.Double(to.Float64(v)))
}

func (n NumberSchema) clamp(value controls.Double) controls.Double {
	if n.Maximum != 0 {
		if value > n.Maximum {
			return n.Maximum
		}
	}
	if value < n.Minimum {
		return n.Minimum
	}
	return value
}

func (n *NumberSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
