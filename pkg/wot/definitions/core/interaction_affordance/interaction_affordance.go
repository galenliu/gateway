package interaction_affordance

import (
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
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

func (v *InteractionAffordance) UnmarshalJSON(data []byte) error {

	node := gjson.GetBytes(data, "@type")
	if node.Exists() {
		v.AtType = node.String()
	}

	node = gjson.GetBytes(data, "title")
	if node.Exists() {
		v.Title = node.String()
	}

	node = gjson.GetBytes(data, "description")
	if node.Exists() {
		v.Description = node.String()
	}

	node = gjson.GetBytes(data, "titles")
	if node.IsArray() {
		if v.Titles == nil {
			v.Titles = make(map[string]string, 0)
		}
		node.ForEach(func(key, value gjson.Result) bool {
			v.Titles[key.String()] = value.String()
			return true
		})
	}

	node = gjson.GetBytes(data, "descriptions")
	if node.IsArray() {
		if v.Descriptions == nil {
			v.Descriptions = make(map[string]string, 0)
		}
		node.ForEach(func(key, value gjson.Result) bool {
			v.Descriptions[key.String()] = value.String()
			return true
		})
	}

	node = gjson.GetBytes(data, "forms")
	if node.IsArray() {
		var f controls.Form
		fdata := data[node.Index : node.Index+len(node.Raw)]
		err := json.Unmarshal(fdata, &f)
		if err != nil {
			v.Forms = append(v.Forms, f)
		}
	}

	node = gjson.GetBytes(data, "uriVariables")
	if node.Exists() && node.IsArray() {
		var s map[string]schema.Schema
		schemaData := data[node.Index : node.Index+len(node.Raw)]
		err := json.Unmarshal(schemaData, &s)
		if err != nil {
			v.UriVariables = s
		}
	}
	return nil
}
