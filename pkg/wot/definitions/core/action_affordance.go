package core

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	"github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ActionAffordance struct {
	*ia.InteractionAffordance
	Input      *data_schema.DataSchema `json:"input,omitempty"`
	Output     *data_schema.DataSchema `json:"output,omitempty"`
	Safe       bool                    `json:"safe,omitempty" wot:"withDefault"`
	Idempotent bool                    `json:"idempotent,omitempty" wot:"withDefault"`
}

type ActionDescription struct {
	AtType       string                            `json:"@type,omitempty,optional"`
	Title        string                            `json:"title,omitempty,optional"`
	Titles       map[string]string                 `json:"titles,omitempty,optional"`
	Description  string                            `json:"description,omitempty,optional"`
	Descriptions map[string]string                 `json:"descriptions,omitempty,optional"`
	Forms        []controls.Form                   `json:"forms,omitempty,mandatory" wot:"optional"`
	UriVariables map[string]data_schema.DataSchema `json:"uriVariables,omitempty"`

	Input      *data_schema.DataSchema `json:"input,omitempty"`
	Output     *data_schema.DataSchema `json:"output,omitempty"`
	Safe       bool                    `json:"safe,omitempty" wot:"withDefault"`
	Idempotent bool                    `json:"idempotent,omitempty" wot:"withDefault"`
}

func FromDescription(desc ActionDescription) ActionAffordance {
	return ActionAffordance{
		InteractionAffordance: &ia.InteractionAffordance{
			AtType:       desc.AtType,
			Title:        desc.Title,
			Titles:       desc.Titles,
			Description:  desc.Description,
			Descriptions: desc.Descriptions,
			Forms:        desc.Forms,
			UriVariables: desc.UriVariables,
		},
		Input:      desc.Input,
		Output:     desc.Output,
		Safe:       desc.Safe,
		Idempotent: desc.Idempotent,
	}
}

func (a *ActionAffordance) MarshalJSON() ([]byte, error) {
	action := ActionDescription{
		AtType:       a.AtType,
		Title:        a.Title,
		Titles:       a.Titles,
		Description:  a.Description,
		Descriptions: a.Descriptions,
		Forms:        a.Forms,
		UriVariables: a.UriVariables,
		Input:        a.Input,
		Output:       a.Output,
		Safe:         a.Safe,
		Idempotent:   a.Idempotent,
	}
	return json.Marshal(action)
}

func (a *ActionAffordance) UnmarshalJSON(data []byte) error {

	var action ActionDescription
	err := json.Unmarshal(data, &action)
	if err != nil {
		return err
	}
	a.InteractionAffordance = &ia.InteractionAffordance{
		AtType:       action.AtType,
		Title:        action.Title,
		Titles:       action.Titles,
		Description:  action.Description,
		Descriptions: action.Descriptions,
		Forms:        action.Forms,
		UriVariables: action.UriVariables,
	}
	a.Input = action.Input
	a.Output = action.Output
	a.Safe = action.Safe
	a.Idempotent = action.Idempotent
	return nil
}
