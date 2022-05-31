package core

import (
	"github.com/bytedance/sonic"
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
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

	err := sonic.Unmarshal(data, a.InteractionAffordance)
	if err != nil {
		return err
	}
	node, _ := sonic.Get(data, "input")
	if node.Exists() {
		d, err := node.MarshalJSON()
		schema, err := data_schema.MarshalSchema(d)
		if err == nil {
			a.Input = schema
		}
	}

	node, _ = sonic.Get(data, "output")
	if node.Exists() {
		d, err := node.MarshalJSON()
		schema, err := data_schema.MarshalSchema(d)
		if err == nil {
			a.Input = schema
		}
	}
	node, _ = sonic.Get(data, "safe")
	if node.Exists() {
		a.Safe, _ = node.Bool()
	}

	node, _ = sonic.Get(data, "idempotent")
	if node.Exists() {
		a.Idempotent, _ = node.Bool()
	}

	return nil
}
