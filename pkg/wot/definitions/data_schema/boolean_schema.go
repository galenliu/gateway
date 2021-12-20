package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type BooleanSchema struct {
	*DataSchema
}

func (schema *BooleanSchema) UnmarshalJSON(data []byte) error {
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	schema.DataSchema = &dataSchema
	return nil
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
