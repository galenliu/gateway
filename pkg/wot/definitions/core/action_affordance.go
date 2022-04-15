package core

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ActionAffordance struct {
	*ia.InteractionAffordance
	Input      data_schema.Schema `json:"input,omitempty"`
	Output     data_schema.Schema `json:"output,omitempty"`
	Safe       bool               `json:"safe,omitempty" wot:"withDefault"`
	Idempotent bool               `json:"idempotent,omitempty" wot:"withDefault"`
}

type ActionSchema = data_schema.Schema

type ActionDescription struct {
	AtType       string                        `json:"@type,omitempty,optional"`
	Title        string                        `json:"title,omitempty,optional"`
	Titles       map[string]string             `json:"titles,omitempty,optional"`
	Description  string                        `json:"description,omitempty,optional"`
	Descriptions map[string]string             `json:"descriptions,omitempty,optional"`
	Forms        []controls.Form               `json:"forms,omitempty,mandatory" wot:"optional"`
	UriVariables map[string]data_schema.Schema `json:"uriVariables,omitempty"`

	Input      ActionSchema       `json:"input,omitempty"`
	Output     data_schema.Schema `json:"output,omitempty"`
	Safe       bool               `json:"safe,omitempty" wot:"withDefault"`
	Idempotent bool               `json:"idempotent,omitempty" wot:"withDefault"`
}

func (a *ActionAffordance) UnmarshalJSON(data []byte) error {
	json.Get(data).ToVal(a.InteractionAffordance)
	a.Input, _ = data_schema.MarshalSchema(json.Get(data, "input"))
	a.Output, _ = data_schema.MarshalSchema(json.Get(data, "input"))
	a.Safe = json.Get(data, "safe").ToBool()
	a.Idempotent = json.Get(data, "idempotent").ToBool()
	return nil
}
