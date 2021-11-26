package core

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ActionAffordance struct {
	*interaction_affordance.InteractionAffordance
	Input      *data_schema.DataSchema `json:"input,omitempty"`
	Output     *data_schema.DataSchema `json:"output,omitempty"`
	Safe       bool                    `json:"safe,omitempty"`       //with default
	Idempotent bool                    `json:"idempotent,omitempty"` //with default
}

func NewActionAffordanceFromString(description string) *ActionAffordance {
	data := []byte(description)
	var a = ActionAffordance{}
	if a.InteractionAffordance = interaction_affordance.NewInteractionAffordanceFromString(description); a.InteractionAffordance == nil {
		return nil
	}

	a.Input = data_schema.NewDataSchemaFromString(json.Get(data, "input").ToString())
	a.Input = data_schema.NewDataSchemaFromString(json.Get(data, "output").ToString())
	a.Safe = json.Get(data, "safe").ToBool()
	a.Idempotent = json.Get(data, "idempotent").ToBool()

	if a.Forms == nil {
		a.Forms = append(a.Forms, controls.Form{
			Href: "",
			Op:   controls.NewArrayOfString(controls.InvokeAction),
		})
	}
	return &a
}
