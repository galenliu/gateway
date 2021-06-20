package data_schema

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type ObjectSchema struct {
	*dataSchema
	Properties map[string]dataSchema `json:"properties,omitempty"`
	Required   []string              `json:"required,omitempty"`
}

func NewObjectSchema() *ObjectSchema {
	obj := &ObjectSchema{}
	obj.Properties = make(map[string]dataSchema)
	obj.dataSchema = &dataSchema{
		Type: hypermedia_controls.Object,
	}
	return obj
}

func NewObjectSchemaFromString(data string) *ObjectSchema {
	var ds dataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	var s = NewObjectSchema()
	m := gjson.Get(data, "properties").Map()
	if len(m) > 0 {
		s.Properties = make(map[string]dataSchema)
		for k, v := range m {
			s.Properties[k] = NewDataSchemaFromString(v.String())
		}
	}
	l := gjson.Get(data, "required").Array()
	if len(l) > 0 {
		for _, d := range l {
			s.Required = append(s.Required, d.String())
		}
	}
	s.dataSchema = &ds
	return s
}

func (n *ObjectSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
