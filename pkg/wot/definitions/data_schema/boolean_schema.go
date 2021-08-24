package data_schema

import controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type BooleanSchema struct {
	*DataSchema
}

func NewBooleanSchemaFromString(description string) *BooleanSchema {
	var schema = BooleanSchema{}
	schema.DataSchema = NewDataSchemaFromString(description)
	if schema.DataSchema == nil || schema.DataSchema.GetType() != controls.TypeString {
		return nil
	}
	return &schema
}

func (b *BooleanSchema) Convert(v interface{}) interface{} {
	return controls.ToBool(v)
}

func (b *BooleanSchema) GetDefaultValue() interface{} {
	if b.DataSchema.Default != nil {
		return b.Convert(b.Default)
	}
	return false
}
