package wot

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type ObjectSchema struct {
	*DataSchema
	Properties map[string]IDataSchema `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

func NewObjectSchema() *ObjectSchema {
	obj := &ObjectSchema{}
	obj.Properties = make(map[string]IDataSchema)
	obj.DataSchema = &DataSchema{
		Type: Object,
	}
	return obj
}

func NewObjectSchemaFromString(data string) *ObjectSchema {
	var ds DataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	var s = NewObjectSchema()
	m := gjson.Get(data, "properties").Map()
	if len(m) > 0 {
		s.Properties = make(map[string]IDataSchema)
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
	s.DataSchema = &ds
	return s
}

func (n *ObjectSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
