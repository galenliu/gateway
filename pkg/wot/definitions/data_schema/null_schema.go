package data_schema

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type NullSchema struct {
	*dataSchema
}

func NewNullSchemaFromString(description string) *NullSchema {
	var schema = NullSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil || schema.dataSchema.GetType() != hypermedia_controls.TypeString {
		return nil
	}
	return &schema
}

func (n NullSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NullSchema) Convert(v interface{}) interface{} {
	return v
}

func (n *NullSchema) GetDefaultValue() interface{} {
	if n.Default != nil {
		return n.Convert(n.Default)
	}
	return nil
}
