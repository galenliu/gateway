package data_schema

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ObjectSchema struct {
	*DataSchema
	Properties map[string]*DataSchema `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

func NewObjectSchemaFromString(description string) *ObjectSchema {
	data := []byte(description)
	var schema = ObjectSchema{}
	schema.DataSchema = NewDataSchemaFromString(description)
	if schema.DataSchema == nil || schema.DataSchema.GetType() != hypermedia_controls.TypeObject {
		return nil
	}
	var properties map[string]string
	json.Get(data, "properties").ToVal(&properties)
	if len(properties) > 0 {
		schema.Properties = make(map[string]*DataSchema)
		for n, p := range properties {
			schema.Properties[n] = NewDataSchemaFromString(p)
		}
	}

	var required []string
	json.Get(data, "required").ToVal(&required)
	if len(required) > 0 {
		for _, r := range required {
			schema.Required = append(schema.Required, r)
		}
	}
	return &schema
}

func (o *ObjectSchema) Convert(v interface{}) interface{} {
	return v
}

func (o *ObjectSchema) GetDefaultValue() interface{} {
	return o.Default
}
