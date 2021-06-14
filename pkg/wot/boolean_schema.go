package wot

import json "github.com/json-iterator/go"

type BooleanSchema struct {
	*DataSchema
}

func NewBooleanSchema() *BooleanSchema {
	b := &BooleanSchema{}
	b.Type = Boolean
	return b
}

func NewBooleanSchemaFromString(data string) *BooleanSchema {
	var ds DataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = BooleanSchema{}
	s.DataSchema = &ds
	return &s
}

func (n BooleanSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
