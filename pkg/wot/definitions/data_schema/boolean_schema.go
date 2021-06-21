package data_schema

import "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type BooleanSchema struct {
	*dataSchema
}

func NewBooleanSchemaFromString(description string) *BooleanSchema {
	var schema = BooleanSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil ||  schema.dataSchema.GetType() != hypermedia_controls.TypeString {
		return nil
	}
	return &schema
}
