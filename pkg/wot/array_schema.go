package wot

import (
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type ArraySchema struct {
	*DataSchema
	Items    []IDataSchema `json:"items,omitempty"`
	MinItems int           `json:"minItems,omitempty"`
	MaxItems int           `json:"maxItems,omitempty"`
}

func NewArraySchemaFromString(data string) *ArraySchema {
	var ds DataSchema
	err := json.Unmarshal([]byte(data), &ds)
	if err != nil {
		return nil
	}
	var s = ArraySchema{}
	arr := gjson.Get(data, "items").Array()
	if len(arr) > 0 {
		for _, r := range arr {
			s.Items = append(s.Items, NewDataSchemaFromString(r.String()))
		}
	}
	s.MinItems = int(gjson.Get(data, "minItems").Int())
	s.MaxItems = int(gjson.Get(data, "maxItems").Int())
	s.DataSchema = &ds
	return &s
}

func NewArraySchema() *ArraySchema {
	d := &ArraySchema{}
	d.Type = Array
	return d
}

func (n *ArraySchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}
