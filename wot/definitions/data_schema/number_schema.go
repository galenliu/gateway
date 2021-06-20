package data_schema

import (
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type NumberSchema struct {
	*dataSchema
	Minimum          float64 `json:"minimum"`
	ExclusiveMinimum float64 `json:"exclusiveMinimum,omitempty"`
	Maximum          float64 `json:"maximum,omitempty"`
	ExclusiveMaximum float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64 `json:"multipleOf,omitempty"`
}

func NewNumberSchemaFromString(data string) *NumberSchema {
	var ds dataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = NewNumberSchema()
	s.Minimum = gjson.Get(data, "minimum").Float()
	s.ExclusiveMinimum = gjson.Get(data, "exclusiveMinimum").Float()
	s.Maximum = gjson.Get(data, "maximum").Float()
	s.ExclusiveMaximum = gjson.Get(data, "exclusiveMaximum").Float()
	s.MultipleOf = gjson.Get(data, "multipleOf").Float()
	s.dataSchema = &ds
	return s
}

func NewNumberSchema() *NumberSchema {
	d := &NumberSchema{}

	return d
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
