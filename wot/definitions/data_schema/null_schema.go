package data_schema

import json "github.com/json-iterator/go"

type NullSchema struct {
	*dataSchema
}

func NewNullSchemaFromString(data string) *NullSchema {
	var ds dataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = NullSchema{}
	s.dataSchema = &ds
	return &s
}

func (n NullSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
