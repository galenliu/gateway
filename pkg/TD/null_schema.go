package TD

import json "github.com/json-iterator/go"

type NullSchema struct {
	*DataSchema
}

func NewNullSchemaFromString(data string) *NullSchema {
	var ds DataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = NullSchema{}
	s.DataSchema = &ds
	return &s
}

func (n NullSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
