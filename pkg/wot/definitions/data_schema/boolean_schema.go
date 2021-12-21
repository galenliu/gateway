package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type BooleanSchema struct {
	*DataSchema
}

func (schema *BooleanSchema) Convert(v any) any {
	return controls.ToBool(v)
}

func (schema *BooleanSchema) GetDefaultValue() any {
	if schema.DataSchema.Default != nil {
		return schema.Convert(schema.Default)
	}
	return false
}
