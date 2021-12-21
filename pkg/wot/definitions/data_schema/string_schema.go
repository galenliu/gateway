package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type StringSchema struct {
	*DataSchema
	MinLength        *controls.UnsignedInt `json:"minLength,omitempty"`
	MaxLength        *controls.UnsignedInt `json:"maxLength,omitempty"`
	Pattern          string                `json:"pattern,omitempty"`
	ContentEncoding  string                `json:"contentEncoding,omitempty"`
	ContentMediaType string                `json:"contentMediaType,omitempty"`
}

func (schema *StringSchema) GetDefaultValue() any {
	if schema.DataSchema.Default != nil {
		return schema.Default
	}
	if len(schema.Enum) > 0 {
		return schema.Enum[0]
	}
	return ""
}

func (schema *StringSchema) verifyType(value any) bool {
	switch value.(type) {
	case string:
		return true
	}
	return false
}

func (schema *StringSchema) clamp(value string) string {
	if schema.MaxLength != nil {
		value = value[0:*schema.MaxLength]
	}
	return value
}
