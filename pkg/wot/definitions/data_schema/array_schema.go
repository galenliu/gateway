package data_schema

import (
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"math"
)

type ArraySchema struct {
	*dataSchema
	Items    []DataSchema         `json:"items,omitempty"`
	MinItems controls.UnsignedInt `json:"minItems,omitempty"`
	MaxItems controls.UnsignedInt `json:"maxItems,omitempty"`
}

func NewArraySchemaFromString(description string) *ArraySchema {
	data := []byte(description)
	var schema = ArraySchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil || schema.dataSchema.GetType() != controls.TypeArray {
		return nil
	}
	var items []string
	json.Get(data, "items").ToVal(&items)
	for _, i := range items {
		schema.Items = append(schema.Items, NewDataSchemaFromString(i))
	}
	schema.MinItems = controls.UnsignedInt(controls.JSONGetUint64(data, "minItems", math.MinInt64))
	schema.MaxItems = controls.UnsignedInt(controls.JSONGetUint64(data, "maxItems", math.MaxUint64))
	return &schema
}

func (a *ArraySchema) Convert(v interface{}) interface{} {
	return v
}

func (a *ArraySchema) GetDefaultValue() interface{} {
	if a.dataSchema.Default != nil {
		return a.Convert(a.Default)
	}
	return nil
}