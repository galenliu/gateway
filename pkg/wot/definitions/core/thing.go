package core

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
	json "github.com/json-iterator/go"
)

type ThingInterface interface {
}

type ThingProperties map[string]PropertyAffordance
type ThingActions map[string]ActionAffordance
type ThingEvents map[string]EventAffordance

type Thing struct {
	Context      controls.URI      `json:"@context,omitempty"` //mandatory
	Title        string            `json:"title,omitempty"`    //mandatory
	Titles       map[string]string `json:"titles,omitempty"`
	Id           controls.URI      `json:"id"`
	Type         []string          `json:"@type"`
	Description  string            `json:"description,omitempty"`
	Descriptions map[string]string `json:"descriptions,omitempty"`

	Support controls.URI `json:"support,omitempty"`
	Base    controls.URI `json:"base,omitempty"`

	Version  *VersionInfo       `json:"version,omitempty"`
	Created  *controls.DataTime `json:"created,omitempty"`
	Modified *controls.DataTime `json:"modified,omitempty"`

	Properties ThingProperties `json:"properties,omitempty"`
	Actions    ThingActions    `json:"actions,omitempty"`
	Events     ThingEvents     `json:"events,omitempty"`

	Links []controls.Link `json:"links,omitempty"`
	Forms []controls.Form `json:"forms,omitempty"`

	Security            controls.ArrayOfString                    `json:"security,omitempty"`            //mandatory
	SecurityDefinitions map[string]*securityScheme.SecurityScheme `json:"securityDefinitions,omitempty"` //mandatory

	Profile           []controls.URI                    `json:"profile,omitempty"`
	SchemaDefinitions map[string]*dataSchema.DataSchema `json:"schemaDefinitions,omitempty"`
}

func NewThingFromString(description string) (thing *Thing) {

	data := []byte(description)
	t := &Thing{}
	t.Id = controls.URI(controls.JSONGetString(data, "id", ""))
	t.Title = controls.JSONGetString(data, "title", "")
	if t.Id == "" {
		t.Id = controls.URI(t.Title)
	}
	t.Context = controls.URI(json.Get(data, "@context").ToString())
	t.Security = controls.NewArrayOfString(json.Get(data, "security").ToString())
	t.Description = controls.JSONGetString(data, "description", "")

	var m map[string]*securityScheme.SecurityScheme
	json.Get(data, "securityDefinitions").ToVal(&m)
	if m != nil {
		t.SecurityDefinitions = m
	}

	if propMap := controls.JSONGetMap(data, "properties"); len(propMap) > 0 {
		t.Properties = make(map[string]PropertyAffordance)
		for name, d := range propMap {
			prop := NewPropertyAffordanceFromString(d)
			t.Properties[name] = prop
		}
		t.Forms = append(t.Forms, controls.Form{Op: controls.ReadallProperties + controls.WriteAllProperties, Href: controls.URI(string(t.Id + constant.PropertiesPath)), ContentType: dataSchema.ApplicationJson})
	}

	if actionMap := controls.JSONGetMap(data, "properties"); len(actionMap) > 0 {
		t.Actions = make(map[string]ActionAffordance)
		for name, data := range actionMap {
			action := NewActionAffordanceFromString(data)
			t.Actions[name] = *action
		}
		t.Forms = append(t.Forms, controls.Form{Href: controls.URI(string(t.Id + constant.ActionsPath))})
	}

	if eventMap := controls.JSONGetMap(data, "events"); len(eventMap) > 0 {
		t.Events = make(map[string]EventAffordance)
		for name, data := range eventMap {
			event := NewEventAffordanceFromString(data)
			t.Events[name] = *event
		}
		t.Forms = append(t.Forms, controls.Form{Href: controls.URI(string(t.Id + constant.ActionsPath))})
	}

	if t.Forms == nil {
		t.Forms = append(t.Forms, controls.Form{ContentType: "text/html", Href: controls.URI(string(t.Id))})
		t.Forms = append(t.Forms, controls.Form{Href: controls.URI(fmt.Sprintf("wss://localhost/%s", t.Id))})
	}
	return t
}

func (t *Thing) GetId() string {
	return t.Id.GetId()
}

func (t *Thing) GetURI() string {
	return t.Id.GetURI()
}

func (t *Thing) GetHref() string {
	return t.Id.GetURI()
}
