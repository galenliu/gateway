package core

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
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

func NewActionAffordanceFromString(description string) *actionAffordance {
	data := []byte(description)
	var a = actionAffordance{}

	if a.InteractionAffordance = NewInteractionAffordanceFromString(description); a.InteractionAffordance == nil {
		return nil
	}

	a.Input = data_schema.NewDataSchemaFromString(json.Get(data, "input").ToString())
	a.Input = data_schema.NewDataSchemaFromString(json.Get(data, "output").ToString())

	if gjson.Get(description, "safe").Exists() {
		s := gjson.Get(description, "safe").Bool()
		a.Safe = s
	}

	if gjson.Get(description, "idempotent").Exists() {
		s := gjson.Get(description, "idempotent").Bool()
		a.Idempotent = s
	}

	if a.Forms == nil {
		a.Forms = append(a.Forms, controls.Form{
			Href: "",
			Op:   []string{controls.InvokeAction},
		})
	}
	return &a
}
