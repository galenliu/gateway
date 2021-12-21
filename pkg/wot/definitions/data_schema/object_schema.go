package data_schema

type ObjectSchema struct {
	*DataSchema
	Properties map[string]DataSchema `json:"properties,omitempty"`
	Required   []string              `json:"required,omitempty"`
}

func (schema *ObjectSchema) Convert(v any) any {
	return v
}

func (schema *ObjectSchema) GetDefaultValue() any {
	return schema.Default
}
