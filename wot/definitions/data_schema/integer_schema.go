package data_schema

import (
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type IntegerSchema struct {
	*DataSchema
	Minimum          int64 `json:"minimum"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
	Maximum          int64 `json:"maximum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       int64 `json:"multipleOf,omitempty"`
}

func NewIntegerSchemaFromString(data string) *IntegerSchema {
	var ds DataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = NewIntegerSchema()
	s.Minimum = gjson.Get(data, "minimum").Int()
	s.ExclusiveMinimum = gjson.Get(data, "exclusiveMinimum").Int()
	s.Maximum = gjson.Get(data, "maximum").Int()
	s.ExclusiveMaximum = gjson.Get(data, "exclusiveMaximum").Int()
	s.MultipleOf = gjson.Get(data, "multipleOf").Int()
	s.DataSchema = &ds
	return s
}

func NewIntegerSchema() *IntegerSchema {
	d := &IntegerSchema{}

	return d
}

func (n IntegerSchema) ClampInt(value int64) int64 {
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

func (n IntegerSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
