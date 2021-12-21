package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type ArraySchema struct {
	*DataSchema
	Items    []DataSchema          `json:"items,omitempty"`
	MinItems *controls.UnsignedInt `json:"minItems,omitempty"`
	MaxItems *controls.UnsignedInt `json:"maxItems,omitempty"`
}

func (schema *ArraySchema) Convert(v any) any {
	return v
}

func (schema *ArraySchema) GetDefaultValue() any {
	if schema.DataSchema.Default != nil {
		return schema.Convert(schema.Default)
	}
	return nil
}
