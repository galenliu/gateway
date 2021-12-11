package data_schema

import (
	"fmt"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type StringSchema struct {
	*DataSchema
	MinLength        *controls.UnsignedInt `json:"minLength,omitempty"`
	MaxLength        *controls.UnsignedInt `json:"maxLength,omitempty"`
	Pattern          string                `json:"pattern,omitempty"`
	ContentEncoding  string                `json:"contentEncoding,omitempty"`
	ContentMediaType string                `json:"contentMediaType,omitempty"`
}

func (schema *StringSchema) UnmarshalJSON(data []byte) error {
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	if &dataSchema == nil {
		return fmt.Errorf("data schema is nil")
	}
	schema.DataSchema = &dataSchema
	if schema.DataSchema == nil || schema.DataSchema.GetType() != controls.TypeString {
		return fmt.Errorf("type must string")
	}

	if min := json.Get(data, "minLength"); min.LastError() == nil {
		m := controls.UnsignedInt(min.ToInt64())
		schema.MinLength = &m
	}

	if max := json.Get(data, "maxLength"); max.LastError() == nil {
		m := controls.UnsignedInt(max.ToInt64())
		schema.MaxLength = &m
	}
	if v := json.Get(data, "contentEncoding"); v.LastError() == nil {
		schema.ContentEncoding = v.ToString()
	}
	if v := json.Get(data, "contentMediaType"); v.LastError() == nil {
		schema.ContentMediaType = v.ToString()
	}
	return nil
}

//func (schema *StringSchema) Convert(v interface{}) interface{} {
//	return schema.clamp(controls.ToString(v))
//}

func (schema *StringSchema) GetDefaultValue() interface{} {
	if schema.DataSchema.Default != nil {
		return schema.Default
	}
	if len(schema.Enum) > 0 {
		return schema.Enum[0]
	}
	return ""
}

//func (schema *StringSchema) clamp(value string) string {
//	if schema.MaxLength != 0 {
//		if schema.MaxLength < controls.ToUnsignedInt(len(value)) {
//			return string([]rune(value)[:schema.MaxLength])
//		}
//	}
//	return value
//}
