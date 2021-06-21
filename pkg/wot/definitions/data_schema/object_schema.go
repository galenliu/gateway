package data_schema

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ObjectSchema struct {
	*dataSchema
	Properties map[string]DataSchema `json:"properties,omitempty"`
	Required   []string              `json:"required,omitempty"`
}

func NewObjectSchemaFromString(description string) *ObjectSchema {
	data := []byte(description)
	var schema = ObjectSchema{}
	schema.dataSchema = newDataSchemaFromString(description)
	if schema.dataSchema == nil ||  schema.dataSchema.GetType() != hypermedia_controls.TypeObject {
		return nil
	}
	var properties map[string]string
	json.Get(data, "properties").ToVal(&properties)
	if len(properties) > 0 {
		schema.Properties = make(map[string]DataSchema)
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

func (n *ObjectSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
