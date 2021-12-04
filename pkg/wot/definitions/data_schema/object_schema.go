package data_schema

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ObjectSchema struct {
	*DataSchema
	Properties map[string]DataSchema `json:"properties,omitempty"`
	Required   []string              `json:"required,omitempty"`
}

func (schema *ObjectSchema) UnmarshalJSON(data []byte) error {
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	if schema.DataSchema == nil || schema.DataSchema.GetType() != hypermedia_controls.TypeObject {
		return fmt.Errorf("type must object")
	}
	var properties map[string]DataSchema
	json.Get(data, "properties").ToVal(&properties)
	if len(properties) > 0 {
		schema.Properties = properties
	}
	var required []string
	json.Get(data, "required").ToVal(&required)
	if len(required) > 0 {
		for _, r := range required {
			schema.Required = append(schema.Required, r)
		}
	}
	return nil
}

func (schema *ObjectSchema) Convert(v interface{}) interface{} {
	return v
}

func (schema *ObjectSchema) GetDefaultValue() interface{} {
	return schema.Default
}
