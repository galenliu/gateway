package interaction_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type InteractionAffordance struct {
	AtType       string                        `json:"@type"`
	Title        string                        `json:"title,omitempty"`
	Titles       map[string]string             `json:"titles,omitempty"`
	Description  string                        `json:"description,omitempty"`
	Descriptions map[string]string             `json:"descriptions,omitempty"`
	Forms        []controls.Form               `json:"forms,omitempty"`
	UriVariables map[string]*schema.DataSchema `json:"uriVariables,omitempty"`
}

func NewInteractionAffordanceFromString(description string) *InteractionAffordance {
	var i = InteractionAffordance{}
	data := []byte(description)

	var uriVariables map[string]*schema.DataSchema
	json.Get(data, "uriVariables").ToVal(&uriVariables)
	if uriVariables != nil {
		i.UriVariables = uriVariables
	}

	i.AtType = controls.JSONGetString(data, "@type", "")
	i.Title = controls.JSONGetString(data, "title", "")
	i.Titles = controls.JSONGetMap(data, "titles")
	i.Description = controls.JSONGetString(data, "description", "")
	i.Descriptions = controls.JSONGetMap(data, "descriptions")

	var forms []controls.Form
	json.Get(data, "forms").ToVal(&forms)
	if len(forms) > 0 {
		i.Forms = forms
	} else {
		return nil
	}
	var uris map[string]*schema.DataSchema
	json.Get(data, "uriVariables").ToVal(&uris)
	if len(uris) > 0 {
		i.UriVariables = uris
	}
	return &i
}
