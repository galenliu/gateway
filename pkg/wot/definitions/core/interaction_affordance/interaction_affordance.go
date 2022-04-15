package interaction_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type InteractionAffordance struct {
	AtType       string                   `json:"@type,omitempty,optional"`
	Title        string                   `json:"title,omitempty,optional"`
	Titles       map[string]string        `json:"titles,omitempty,optional"`
	Description  string                   `json:"description,omitempty,optional"`
	Descriptions map[string]string        `json:"descriptions,omitempty,optional"`
	Forms        []controls.Form          `json:"forms,omitempty,mandatory" wot:"optional"`
	UriVariables map[string]schema.Schema `json:"uriVariables,omitempty"`
}

func (v InteractionAffordance) UnmarshalJSON(data []byte) error {
	v.AtType = json.Get(data, "@type").ToString()
	v.Title = json.Get(data, "title").ToString()
	json.Get(data, "titles").ToVal(&v.Titles)
	v.Description = json.Get(data, "description").ToString()
	json.Get(data, "descriptions").ToVal(&v.Descriptions)
	json.Get(data, "forms").ToVal(&v.Forms)
	var uriVariablesMap map[string]json.Any
	var uriVariables = make(map[string]schema.Schema, 0)
	json.Get(data, "uriVariables").ToVal(&uriVariablesMap)
	for n, u := range uriVariablesMap {
		s, e := schema.MarshalSchema(u)
		if e != nil {
			uriVariables[n] = s
		}
	}
	v.UriVariables = uriVariables
	return nil
}
