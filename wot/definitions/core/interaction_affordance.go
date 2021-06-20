package core

import (
	"github.com/galenliu/gateway/wot/definitions/data_schema"
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"log"
)

type InteractionAffordance struct {
	AtType       string                                     `json:"@type"`
	Title        string                                     `json:"title,omitempty"`
	Titles       map[string]string                          `json:"titles,omitempty"`
	Description  string                                     `json:"description,omitempty"`
	Descriptions map[string]string                          `json:"descriptions,omitempty"`
	Forms        []hypermedia_controls.Form                 `json:"forms,omitempty"`
	UriVariables map[string]data_schema.DataSchemaInterface `json:"uriVariables,omitempty"`
}

func NewInteractionAffordanceFromString(description string) *InteractionAffordance {
	var i = InteractionAffordance{}

	if gjson.Get(description, "uriVariables").Exists() {
		m := gjson.Get(description, "uriVariables").Map()
		if len(m) > 0 {
			i.UriVariables = make(map[string]data_schema.DataSchemaInterface)
			for k, v := range m {
				i.UriVariables[k] = data_schema.NewDataSchemaFromString(v.String())
			}
		}
	}
	if gjson.Get(description, "title").Exists() {
		i.Title = gjson.Get(description, "title").String()
	}
	if gjson.Get(description, "@type").Exists() {
		i.AtType = gjson.Get(description, "@type").String()
	}

	if gjson.Get(description, "titles").Exists() {
		for k, v := range gjson.Get(description, "title").Map() {
			i.Titles[k] = v.String()
		}
	}

	if gjson.Get(description, "descriptions").Exists() {
		for k, v := range gjson.Get(description, "descriptions").Map() {
			i.Descriptions[k] = v.String()
		}
	}

	if gjson.Get(description, "forms").Exists() {
		for _, v := range gjson.Get(description, "forms").Array() {
			var f hypermedia_controls.Form
			err := json.Unmarshal([]byte(v.String()), &f)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			i.Forms = append(i.Forms, f)
		}
	}

	return &i
}

func NewInteractionAffordance() *InteractionAffordance {
	ia := &InteractionAffordance{}
	return ia
}

func (i *InteractionAffordance) MarshalJSON() ([]byte, error) {
	return json.Marshal(i)
}
