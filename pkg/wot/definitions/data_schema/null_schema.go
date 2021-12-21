package data_schema

type NullSchema struct {
	*DataSchema
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
