package core

import (
	wot_properties "github.com/galenliu/gateway/pkg/wot/definitions/core/property_affordance"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
	json "github.com/json-iterator/go"
)

type ThingInterface any

type ThingProperties map[string]PropertyAffordance
type ThingActions map[string]*ActionAffordance
type ThingEvents map[string]*EventAffordance
type ThingSecurityDefinitions map[string]securityScheme.SecurityScheme
type ThingSchemaDefinitions map[string]dataSchema.DataSchema

type Thing struct {
	AtContext    controls.URI           `json:"@context,omitempty" wot:"mandatory"`
	Title        string                 `json:"title,omitempty" wot:"mandatory"`
	Titles       controls.MultiLanguage `json:"titles,omitempty" wot:"optional"`
	Id           controls.URI           `json:"id" wot:"optional"`
	AtType       controls.ArrayOrString `json:"@type"`
	Description  string                 `json:"description,omitempty" wot:"optional"`
	Descriptions controls.MultiLanguage `json:"descriptions,omitempty" wot:"optional"`

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

	Security controls.ArrayOrString `json:"security,omitempty" wot:"mandatory"`

	Profile           []controls.URI         `json:"profile,omitempty" wot:"optional"`
	SchemaDefinitions ThingSchemaDefinitions `json:"schemaDefinitions,omitempty" wot:"optional"`
}

func (props *ThingProperties) UnmarshalJSON(data []byte) error {
	propsAff := make(ThingProperties)
	mapProperties := json.Get(data)
	for _, name := range json.Get(data).Keys() {
		prop := mapProperties.Get(name)
		typ := mapProperties.Get(name, "type").ToString()
		if prop.LastError() != nil || typ == "" {
			continue
		}
		switch typ {
		case controls.TypeString:
			var p wot_properties.StringPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		case controls.TypeBoolean:
			var p wot_properties.BooleanPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		case controls.TypeInteger:
			var p wot_properties.IntegerPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		case controls.TypeNumber:
			var p wot_properties.NumberPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		case controls.TypeObject:
			var p wot_properties.ObjectPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		case controls.TypeArray:
			var p wot_properties.ArrayPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		case controls.TypeNull:
			var p wot_properties.NullPropertyAffordance
			prop.ToVal(&p)
			if &p != nil {
				propsAff[name] = &p
			}
		default:
			continue
		}
	}
	if props != nil && len(propsAff) > 0 {
		*props = propsAff
	}
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

func (t *Thing) GetProperty(name string) PropertyAffordance {
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
