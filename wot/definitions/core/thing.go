package core

import (
	"github.com/galenliu/gateway/wot/definitions/data_schema"
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
	"strings"
	"time"
)

type ThingInterface interface {
}

type Thing struct {
	AtContext    []string `json:"@context"`
	Title        string   `json:"title"`
	Titles       []string `json:"titles,omitempty"`
	ID           string   `json:"id"`
	AtType       []string `json:"@type"`
	Description  string   `json:"description,omitempty"`
	Descriptions []string `json:"descriptions,omitempty"`

	Forms []hypermedia_controls.Form `json:"forms,omitempty"`
	Links []hypermedia_controls.Link `json:"links,omitempty"`

	Support interface{} `json:"support,omitempty"`
	Base    interface{} `json:"base"`

	Version  VersionInfo `json:"version,omitempty"`
	Created  *time.Time  `json:"created,omitempty"`
	Modified *time.Time  `json:"modified,omitempty"`

	Properties map[string]PropertyAffordance `json:"properties,omitempty"`
	Actions    map[string]ActionAffordance   `json:"_actions,omitempty"`
	Events     map[string]EventAffordance    `json:"events,omitempty"`

	Security []string `json:"security,omitempty"`

	SchemaDefinitions []data_schema.dataSchema
}

func (t *Thing) GetID() string {
	sl := strings.Split(t.ID, "/")
	tid := sl[len(sl)-1]
	return tid
}

func (t *Thing) GetThingID() string {
	return t.ID
}
