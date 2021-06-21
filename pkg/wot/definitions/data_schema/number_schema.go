package data_schema

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type NumberSchema struct {
	*dataSchema
	Minimum          float64 `json:"minimum"`
	ExclusiveMinimum float64 `json:"exclusiveMinimum,omitempty"`
	Maximum          float64 `json:"maximum,omitempty"`
	ExclusiveMaximum float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64 `json:"multipleOf,omitempty"`
}

func NewNumberSchemaFromString(description string) *NumberSchema {
	data := []byte(description)
	var schema = NumberSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil ||  schema.dataSchema.GetType() != hypermedia_controls.TypeNumber {
		return nil
	}

	schema.Minimum = json.Get(data, "minimum").ToFloat64()
	schema.ExclusiveMinimum = json.Get(data, "exclusiveMinimum").ToFloat64()
	schema.Maximum = json.Get(data, "maximum").ToFloat64()
	schema.ExclusiveMaximum = json.Get(data, "exclusiveMaximum").ToFloat64()
	schema.MultipleOf = json.Get(data, "multipleOf").ToFloat64()
	return &schema
}

func (n NumberSchema) ClampFloat(value float64) float64 {
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
