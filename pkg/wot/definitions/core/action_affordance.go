package core

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	json "github.com/json-iterator/go"
)

type ActionAffordance struct {
	*ia.InteractionAffordance
	Input      *data_schema.DataSchema `json:"input,omitempty"`
	Output     *data_schema.DataSchema `json:"output,omitempty"`
	Safe       bool                    `json:"safe,omitempty"`       //with default
	Idempotent bool                    `json:"idempotent,omitempty"` //with default
}

func (a *ActionAffordance) UnmarshalJSON(data []byte) error {
	var i ia.InteractionAffordance
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	var input data_schema.DataSchema
	json.Get(data, "input").ToVal(&input)
	a.Input = &input

	var output data_schema.DataSchema
	json.Get(data, "output").ToVal(&output)
	a.Output = &output

	a.Safe = json.Get(data, "safe").ToBool()
	a.Idempotent = json.Get(data, "idempotent").ToBool()
	return nil
}
