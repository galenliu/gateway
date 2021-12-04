package core

import (
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
type ThingSecurityDefinitions map[string]securityScheme.SecurityScheme
type ThingSchemaDefinitions map[string]dataSchema.DataSchema

type Thing struct {
	Context      controls.URI      `json:"@context,omitempty" wot:"mandatory"`
	Title        string            `json:"title,omitempty" wot:"mandatory"`
	Titles       map[string]string `json:"titles,omitempty" wot:"optional"`
	Id           controls.URI      `json:"id" wot:"optional"`
	Type         []string          `json:"@type"`
	Description  string            `json:"description,omitempty" wot:"optional"`
	Descriptions map[string]string `json:"descriptions,omitempty" wot:"optional"`

	Support controls.URI `json:"support,omitempty" wot:"optional"`
	Base    controls.URI `json:"base,omitempty" wot:"optional"`

	Version  *VersionInfo       `json:"version,omitempty" wot:"optional"`
	Created  *controls.DataTime `json:"created,omitempty" wot:"optional"`
	Modified *controls.DataTime `json:"modified,omitempty" wot:"optional"`

	Properties ThingProperties `json:"properties,omitempty" wot:"optional"`
	Actions    ThingActions    `json:"actions,omitempty" wot:"optional"`
	Events     ThingEvents     `json:"events,omitempty" wot:"optional"`

	Links []controls.Link `json:"links,omitempty"`
	Forms []controls.Form `json:"forms,omitempty"`

	Security            controls.ArrayOfString   `json:"security,omitempty" wot:"mandatory"`
	SecurityDefinitions ThingSecurityDefinitions `json:"securityDefinitions,omitempty" wot:"mandatory"`

	Profile           []controls.URI         `json:"profile,omitempty" wot:"optional"`
	SchemaDefinitions ThingSchemaDefinitions `json:"schemaDefinitions,omitempty" wot:"optional"`
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	t.Id = controls.URI(json.Get(data, "id").ToString())
	t.Title = json.Get(data, "title").ToString()
	if t.Title == "" {
		return nil
	}
	if t.Id == "" {
		t.Id = controls.URI(t.Title)
	}
	t.Context = controls.URI(json.Get(data, "@context").ToString())
	t.Security = controls.NewArrayOfString(json.Get(data, "security").ToString())
	t.Description = json.Get(data, "description").ToString()

	var m map[string]securityScheme.SecurityScheme
	json.Get(data, "securityDefinitions").ToVal(&m)
	if m != nil {
		t.SecurityDefinitions = m
	}

	var props ThingProperties
	json.Get(data, "properties").ToVal(&props)
	t.Properties = props

	var actions ThingActions
	json.Get(data, "actions").ToVal(&actions)
	t.Actions = actions

	var events ThingEvents
	json.Get(data, "events").ToVal(&events)
	t.Events = events

	var links []controls.Link
	json.Get(data, "links").ToVal(&links)
	t.Links = links

	var forms []controls.Form
	json.Get(data, "forms").ToVal(&forms)
	t.Forms = forms

	var security controls.ArrayOfString
	json.Get(data, "security").ToVal(&security)
	t.Security = security

	var securityDefinitions ThingSecurityDefinitions
	json.Get(data, "securityDefinitions").ToVal(&securityDefinitions)
	t.SecurityDefinitions = securityDefinitions

	var profile []controls.URI
	json.Get(data, "profile").ToVal(&profile)
	t.Profile = profile

	var schemaDefinitions ThingSchemaDefinitions
	json.Get(data, "schemaDefinitions").ToVal(&schemaDefinitions)
	t.SchemaDefinitions = schemaDefinitions

	return nil
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
