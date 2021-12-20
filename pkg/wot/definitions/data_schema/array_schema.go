package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ArraySchema struct {
	*DataSchema
	Items    []DataSchema          `json:"items,omitempty"`
	MinItems *controls.UnsignedInt `json:"minItems,omitempty"`
	MaxItems *controls.UnsignedInt `json:"maxItems,omitempty"`
}

func (schema *ArraySchema) UnmarshalJSON(data []byte) error {
	var dataSchema DataSchema
	err := json.Unmarshal(data, &dataSchema)
	if err != nil {
		return err
	}
	schema.DataSchema = &dataSchema
	if schema.DataSchema.GetType() != controls.TypeArray {
		return nil
	}
	var items []DataSchema
	json.Get(data, "items").ToVal(&items)
	if items != nil || len(items) > 0 {
		schema.Items = items
	}
	if minItems := json.Get(data, "minItems"); minItems.ValueType() == json.NumberValue {
		var min = controls.UnsignedInt(minItems.ToInt64())
		schema.MinItems = &min
	}
	if minItems := json.Get(data, "maxItems"); minItems.ValueType() == json.NumberValue {
		var max = controls.UnsignedInt(minItems.ToInt64())
		schema.MaxItems = &max
	}
	return nil
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
