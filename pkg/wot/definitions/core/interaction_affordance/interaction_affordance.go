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
	Forms        []controls.Form              `json:"forms,omitempty,mandatory" wot:"optional"`
	UriVariables map[string]schema.DataSchema `json:"uriVariables,omitempty"`
}

func (i *InteractionAffordance) UnmarshalJSON(data []byte) error {
	var uriVariables map[string]schema.DataSchema
	if uriVar := json.Get(data, "uriVariables"); uriVar.LastError() == nil {
		uriVar.ToVal(&uriVariables)
		if &uriVar != nil && len(uriVariables) > 0 {
			i.UriVariables = uriVariables
		}
	}
	i.Type = json.Get(data, "@type").ToString()
	i.Title = json.Get(data, "title").ToString()
	i.Titles = controls.JSONGetMap(data, "titles")
	i.Description = json.Get(data, "description").ToString()
	i.Descriptions = controls.JSONGetMap(data, "descriptions")

	var forms []controls.Form
	if f := json.Get(data, "forms"); f.LastError() == nil {
		f.ToVal(&forms)
		if &f != nil && len(forms) > 0 {
			i.Forms = forms
		}
	}
	return nil
}
