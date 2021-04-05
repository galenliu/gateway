package models

import (
	"addon/wot"
)

type Property struct {
	*wot.DataSchema
	Name string `json:"name"`

	ReadOnly bool `json:"readOnly"`
	Visible  bool `json:"visible"`

	Minimum interface{} `json:"minimum,omitempty"`
	Maximum interface{} `json:"maximum,omitempty"`

	ThingId string `json:"thingId"`
}
