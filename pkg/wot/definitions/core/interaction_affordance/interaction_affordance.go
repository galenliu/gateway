package interaction_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type InteractionAffordance struct {
	Type         string                       `json:"@type,omitempty,optional"`
	Title        string                       `json:"title,omitempty,optional"`
	Titles       map[string]string            `json:"titles,omitempty,optional"`
	Description  string                       `json:"description,omitempty,optional"`
	Descriptions map[string]string            `json:"descriptions,omitempty,optional"`
	Forms        []controls.Form              `json:"forms,omitempty,mandatory"`
	UriVariables map[string]schema.DataSchema `json:"uriVariables,omitempty"`
}

func (i *InteractionAffordance) UnmarshalJSON(data []byte) error {
	var uriVariables map[string]schema.DataSchema
	json.Get(data, "uriVariables").ToVal(&uriVariables)
	if uriVariables != nil {
		i.UriVariables = uriVariables
	}
	i.Type = json.Get(data, "@type").ToString()
	i.Title = json.Get(data, "@title").ToString()
	i.Titles = controls.JSONGetMap(data, "titles")
	i.Description = json.Get(data, "@description").ToString()
	i.Descriptions = controls.JSONGetMap(data, "descriptions")

	var forms []controls.Form
	json.Get(data, "forms").ToVal(&forms)
	if len(forms) > 0 {
		i.Forms = forms
	} else {
		return nil
	}
	var uris map[string]schema.DataSchema
	json.Get(data, "uriVariables").ToVal(&uris)
	if len(uris) > 0 {
		i.UriVariables = uris
	}
	return nil
}
