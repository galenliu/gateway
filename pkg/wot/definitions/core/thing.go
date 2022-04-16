package core

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core/property_affordance"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
)

type ThingProperties map[string]*property_affordance.PropertyAffordance
type ThingActions map[string]*ActionAffordance
type ThingEvents map[string]*EventAffordance
type ThingSecurityDefinitions map[string]securityScheme.SecurityScheme
type ThingSchemaDefinitions map[string]dataSchema.DataSchema

type Thing struct {
	AtContext controls.URI           `json:"@context,omitempty" wot:"mandatory"`
	AtType    controls.ArrayOrString `json:"@type" wot:"optional"`
	Id        controls.URI           `json:"id" wot:"optional"`
	Title     string                 `json:"title,omitempty" wot:"mandatory"`
	Titles    controls.MultiLanguage `json:"titles,omitempty" wot:"optional"`

	Description  string                 `json:"description,omitempty" wot:"optional"`
	Descriptions controls.MultiLanguage `json:"descriptions,omitempty" wot:"optional"`

	Support controls.URI `json:"support,omitempty" wot:"optional"`
	Base    controls.URI `json:"base,omitempty" wot:"optional"`

	Version  *VersionInfo       `json:"version,omitempty" wot:"optional"`
	Created  *controls.DataTime `json:"created,omitempty" wot:"optional"`
	Modified *controls.DataTime `json:"-" wot:"optional"`

	Properties ThingProperties `json:"properties,omitempty" wot:"optional"`
	Actions    ThingActions    `json:"actions,omitempty" wot:"optional"`
	Events     ThingEvents     `json:"events,omitempty" wot:"optional"`

	Links []controls.Link `json:"links,omitempty"`
	Forms []controls.Form `json:"forms,omitempty"`

	Security            controls.ArrayOrString   `json:"security,omitempty" wot:"mandatory"`
	SecurityDefinitions ThingSecurityDefinitions `json:"securityDefinitions,omitempty" wot:"optional"`

	Profile           []controls.URI         `json:"profile,omitempty" wot:"optional"`
	SchemaDefinitions ThingSchemaDefinitions `json:"schemaDefinitions,omitempty" wot:"optional"`
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

func (t *Thing) GetProperty(name string) *property_affordance.PropertyAffordance {
	p, ok := t.Properties[name]
	if !ok {
		return nil
	}
	return p
}

func (t *Thing) GetAction(name string) *ActionAffordance {
	action, ok := t.Actions[name]
	if ok {
		return action
	}
	return nil
}
