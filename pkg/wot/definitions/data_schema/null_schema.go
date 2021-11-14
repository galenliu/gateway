package data_schema

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type NullSchema struct {
	*DataSchema
}

func NewNullSchemaFromString(description string) *NullSchema {
	var schema = NullSchema{}
	schema.DataSchema = NewDataSchemaFromString(description)
	if schema.DataSchema == nil || schema.DataSchema.GetType() != hypermedia_controls.TypeString {
		return nil
	}
	return &schema
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
