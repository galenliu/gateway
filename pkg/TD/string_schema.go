package TD

import (
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type StringSchema struct {
	*DataSchema
	MinLength int64 `json:"minLength,omitempty"`
	MaxLength int64 `json:"maxLength,omitempty"`
}

func NewStringSchemaFromString(data string) *StringSchema {
	var ds DataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = NewStringSchema()
	s.MinLength = gjson.Get(data, "minLength").Int()
	s.MaxLength = gjson.Get(data, "maxLength").Int()

	s.DataSchema = &ds
	return s
}

func NewStringSchema() *StringSchema {
	d := &StringSchema{}
	return d
}

func (s *StringSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}
