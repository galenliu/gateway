package wot

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type ActionAffordance struct {
	*InteractionAffordance

	Input      IDataSchema `json:"input,omitempty"`
	Output     IDataSchema `json:"output,omitempty"`
	Safe       bool        `json:"safe,omitempty"`
	Idempotent bool        `json:"idempotent,omitempty"`
}

func NewActionAffordanceFromString(data string) *ActionAffordance {
	var ia = InteractionAffordance{}
	err := json.Unmarshal([]byte(data), &ia)
	if err != nil {
		fmt.Print(err.Error())
		return nil
	}
	var a = ActionAffordance{}

	if gjson.Get(data, "input").Exists() {
		s := gjson.Get(data, "input").String()
		d := NewDataSchemaFromString(s)
		if d != nil {
			a.Input = d
		}
	}

	if gjson.Get(data, "output").Exists() {
		s := gjson.Get(data, "output").String()
		d := NewDataSchemaFromString(s)
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
	return &a
}

func NewActionAffordance() *ActionAffordance {
	aa := &ActionAffordance{InteractionAffordance: NewInteractionAffordance()}
	return aa
}
