package core

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	hypermedia_controls2 "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type ActionAffordance interface {
}

type actionAffordance struct {
	*InteractionAffordance

	Input      data_schema.DataSchema `json:"input,omitempty"`
	Output     data_schema.DataSchema `json:"output,omitempty"`
	Safe       bool                   `json:"safe,omitempty"`
	Idempotent bool                   `json:"idempotent,omitempty"`
}

func NewActionAffordanceFromString(data string) ActionAffordance {
	var ia = InteractionAffordance{}
	err := json.Unmarshal([]byte(data), &ia)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	var a = actionAffordance{}

	if gjson.Get(data, "input").Exists() {
		s := gjson.Get(data, "input").String()
		d := data_schema.NewDataSchemaFromString(s)
		if d != nil {
			a.Input = d
		}
	}

	if gjson.Get(data, "output").Exists() {
		s := gjson.Get(data, "output").String()
		d := data_schema.NewDataSchemaFromString(s)
		if d != nil {
			a.Output = d
		}
	}

	if gjson.Get(data, "safe").Exists() {
		s := gjson.Get(data, "safe").Bool()
		a.Safe = s
	}

	if gjson.Get(data, "idempotent").Exists() {
		s := gjson.Get(data, "idempotent").Bool()
		a.Idempotent = s
	}

	if a.Forms == nil {
		a.Forms = append(a.Forms, hypermedia_controls2.Form{
			Href: "",
			Op:   []string{hypermedia_controls2.InvokeAction},
		})
	}
	return &a
}

func NewActionAffordance() *actionAffordance {
	aa := &actionAffordance{InteractionAffordance: NewInteractionAffordance()}
	return aa
}
