package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type StringSchema struct {
	*dataSchema
	MinLength        controls.UnsignedInt `json:"minLength,omitempty"`
	MaxLength        controls.UnsignedInt `json:"maxLength,omitempty"`
	Pattern          string               `json:"pattern,omitempty"`
	ContentEncoding  string               `json:"contentEncoding"`
	ContentMediaType string               `json:"contentMediaType"`
}

func NewStringSchemaFromString(description string) *StringSchema {
	data := []byte(description)
	var schema = StringSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil || schema.dataSchema.GetType() != controls.TypeString {
		return nil
	}
	if json.Get(data, "minLength").ValueType() != json.NilValue {
		schema.MinLength = controls.UnsignedInt(json.Get(data, "minLength").ToInt64())
	}
	if json.Get(data, "maxLength").ValueType() != json.NilValue {
		schema.MinLength = controls.UnsignedInt(json.Get(data, "maxLength").ToInt64())
	}

	if json.Get(data, "contentEncoding").ValueType() != json.StringValue {
		schema.ContentEncoding = json.Get(data, "maxLength").ToString()
	}

	if json.Get(data, "contentMediaType").ValueType() != json.StringValue {
		schema.ContentMediaType = json.Get(data, "maxLength").ToString()
	}
	return &schema
}

func (s *StringSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StringSchema) Convert(v interface{}) interface{} {
	return s.clamp(controls.ToString(v))
}

func (s *StringSchema) GetDefaultValue() interface{} {
	if s.dataSchema.Default != nil {
		return s.Convert(s.dataSchema.Default)
	}
	return ""
}

func (s *StringSchema) clamp(value string) string {
	if s.MaxLength != 0 {
		if s.MaxLength < controls.ToUnsignedInt(len(value)) {
			return string([]rune(value)[:s.MaxLength])
		}
	}
	return value
}
