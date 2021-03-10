package thing

import "gateway/pkg/util"

type Property struct {
	Name        string `json:"name"`
	AtType      string `json:"@type"`
	Type        string `json:"type"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Unit     string `json:"unit,omitempty"`
	ReadOnly bool   `json:"readOnly"`
	Visible  bool   `json:"visible"`

	Minimum interface{} `json:"minimum,omitempty"`
	Maximum interface{} `json:"maximum,omitempty"`
	Value   interface{} `json:"-"`
	Enum    []string    `json:"enum,omitempty"`

	Forms []util.Form `json:"forms"`

	ThingId string `json:"thingId"`
}
