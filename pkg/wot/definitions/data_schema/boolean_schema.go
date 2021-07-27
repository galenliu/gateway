package data_schema

import controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type BooleanSchema struct {
	*dataSchema
}

func NewBooleanSchemaFromString(description string) *BooleanSchema {
	var schema = BooleanSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil || schema.dataSchema.GetType() != controls.TypeString {
		return nil
	}
	return &schema
}

func (b *BooleanSchema) Convert(v interface{}) interface{} {
	return controls.ToBool(v)
}

func (b *BooleanSchema) GetDefaultValue() interface{} {
	if b.dataSchema.Default != nil {
		return b.Convert(b.Default)
	}
	return false
}
