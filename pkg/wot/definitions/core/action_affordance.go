package core

import (
	"encoding/json"
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/tidwall/gjson"
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

	err := json.Unmarshal(data, a.InteractionAffordance)
	if err != nil {
		return err
	}
	node := gjson.GetBytes(data, "input")
	if node.Exists() {
		inputData := data[node.Index : node.Index+len(node.Raw)]
		schema, err := data_schema.UnmarshalSchema(inputData)
		if err == nil {
			a.Input = schema
		}
	}

	node = gjson.GetBytes(data, "output")
	if node.Exists() {
		outputData := data[node.Index : node.Index+len(node.Raw)]
		schema, err := data_schema.UnmarshalSchema(outputData)
		if err == nil {
			a.Input = schema
		}
	}

	node = gjson.GetBytes(data, "safe")
	if node.Exists() {
		a.Safe = node.Bool()
	}

	node = gjson.GetBytes(data, "idempotent")
	if node.Exists() {
		a.Idempotent = node.Bool()
	}
	return nil
}
