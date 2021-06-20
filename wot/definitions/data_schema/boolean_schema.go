package data_schema

import (
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type BooleanSchema struct {
	*dataSchema
}

func NewBooleanSchema() *BooleanSchema {
	b := &BooleanSchema{}
	b.Type = hypermedia_controls.Boolean
	return b
}

func NewBooleanSchemaFromString(data string) *BooleanSchema {
	var ds dataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = BooleanSchema{}
	s.dataSchema = &ds
	return &s
}

func (n BooleanSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
