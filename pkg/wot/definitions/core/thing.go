package core

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/util"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
	json "github.com/json-iterator/go"
)

type ThingInterface interface {
}

type Thing struct {
	AtContext    []string          `json:"@context"`
	Title        string            `json:"title"`
	Titles       map[string]string `json:"titles,omitempty"`
	ID           controls.URI      `json:"id"`
	AtType       []string          `json:"@type"`
	Description  string            `json:"description,omitempty"`
	Descriptions map[string]string `json:"descriptions,omitempty"`

	Forms []controls.Form `json:"forms,omitempty"`
	Links []controls.Link `json:"links,omitempty"`

	Support controls.URI `json:"support,omitempty"`
	Base    controls.URI `json:"base"`

	Version  VersionInfo       `json:"version,omitempty"`
	Created  controls.DataTime `json:"created,omitempty"`
	Modified controls.DataTime `json:"modified,omitempty"`

	Properties map[string]PropertyAffordance `json:"properties,omitempty"`
	Actions    map[string]ActionAffordance   `json:"actions,omitempty"`
	Events     map[string]EventAffordance    `json:"events,omitempty"`

	Security            []string                                 `json:"security,omitempty"`
	SecurityDefinitions map[string]securityScheme.SecurityScheme `json:"securityDefinitions"`
	SchemaDefinitions   []dataSchema.DataSchema
}

func NewThingFromString(description string) (thing *Thing, err error) {

	data := []byte(description)
	t := &Thing{}

	t.ID = controls.URI(controls.JSONGetString(data, "id", ""))
	t.Title = controls.JSONGetString(data, "title", "")
	if t.ID == "" {
		t.ID = controls.URI(t.Title)
	}
	t.AtContext = controls.JSONGetArray(data, "@context")
	t.Security = controls.JSONGetArray(data, "security")
	t.Description = controls.JSONGetString(data, "description", "")

	var m map[string]securityScheme.SecurityScheme
	json.Get(data, "securityDefinitions").ToVal(&m)
	if &m != nil {
		t.SecurityDefinitions = m
	}

	if propMap := controls.JSONGetMap(data, "properties"); len(propMap) > 0 {
		t.Properties = make(map[string]PropertyAffordance)
		for name, d := range propMap {
			prop := NewPropertyAffordanceFromString(d)
			t.Properties[name] = prop
		}
		t.Forms = append(t.Forms, controls.Form{Op: []string{controls.ReadallProperties, controls.WriteAllProperties}, Href: controls.URI(string(t.ID + util.PropertiesPath)), ContentType: dataSchema.ApplicationJson})
	}

	if actionMap := controls.JSONGetMap(data, "properties"); len(actionMap) > 0 {
		t.Actions = make(map[string]ActionAffordance)
		for name, data := range actionMap {
			action := NewActionAffordanceFromString(data)
			t.Actions[name] = action
		}
		t.Forms = append(t.Forms, controls.Form{Href: controls.URI(string(t.ID + util.ActionsPath))})
	}

	if eventMap := controls.JSONGetMap(data, "events"); len(eventMap) > 0 {
		t.Events = make(map[string]EventAffordance)
		for name, data := range eventMap {
			event := NewEventAffordanceFromString(data)
			t.Events[name] = event
		}
		t.Forms = append(t.Forms, controls.Form{Href: controls.URI(string(t.ID + util.ActionsPath))})
	}

	if t.Forms == nil {
		t.Forms = append(t.Forms, controls.Form{ContentType: "text/html", Href: controls.URI(string(t.ID))})
		t.Forms = append(t.Forms, controls.Form{Href: controls.URI(fmt.Sprintf("wss://localhost/%s", t.ID))})
	}
	return t, nil
}

func (t *Thing) GetID() string {
	return t.ID.GetID()
}

func (t *Thing) GetURI() string {
	return t.ID.GetURI()
}
