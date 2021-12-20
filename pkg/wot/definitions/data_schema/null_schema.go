package data_schema

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type NullSchema struct {
	*DataSchema
}

func (schema *NullSchema) UnmarshalJSON(data []byte) error {
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	schema.DataSchema = &dataSchema
	if schema.DataSchema == nil || schema.DataSchema.GetType() != hypermedia_controls.TypeNull {
		return fmt.Errorf("type must null")
	}
	return nil
}

func (schema *NullSchema) Convert(v any) any {
	return v
}

func (schema *NullSchema) GetDefaultValue() any {
	if schema.Default != nil {
		return schema.Convert(schema.Default)
	}
	return nil
}
