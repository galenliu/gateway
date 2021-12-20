package core

import (
	"fmt"
	wot_properties "github.com/galenliu/gateway/pkg/wot/definitions/core/property_affordance"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
	json "github.com/json-iterator/go"
)

type ThingInterface any

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

	Version  *VersionInfo      `json:"version,omitempty" wot:"optional"`
	Created  controls.DataTime `json:"created,omitempty" wot:"optional"`
	Modified controls.DataTime `json:"modified,omitempty" wot:"optional"`

	Properties ThingProperties `json:"properties,omitempty" wot:"optional"`
	Actions    ThingActions    `json:"actions,omitempty" wot:"optional"`
	Events     ThingEvents     `json:"events,omitempty" wot:"optional"`

	Links []controls.Link `json:"links,omitempty"`
	Forms []controls.Form `json:"forms,omitempty"`

	Security            controls.ArrayOrString   `json:"security,omitempty" wot:"mandatory"`
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
	t.Security = controls.ArrayOrString(json.Get(data, "security").ToString())
	t.Description = json.Get(data, "description").ToString()
	t.Support = controls.URI(json.Get(data, "support").ToString())

	var typ []string
	if tp := json.Get(data, "@type"); tp.LastError() == nil {
		tp.ToVal(&typ)
		if tp != nil {
			t.Type = typ
		}
	} else {
		return fmt.Errorf("type nil")
	}

	var m map[string]securityScheme.SecurityScheme
	if s := json.Get(data, "securityDefinitions"); s.LastError() == nil {
		s.ToVal(&m)
		if m != nil {
			t.SecurityDefinitions = m
		}
	}

	var props ThingProperties
	if properties := json.Get(data, "properties"); properties.LastError() == nil {
		props = make(ThingProperties)
		for _, name := range properties.Keys() {
			prop := properties.Get(name)
			typ := properties.Get(name, "type").ToString()
			if prop.LastError() != nil || typ == "" {
				continue
			}
			switch typ {
			case controls.TypeString:
				var p wot_properties.StringPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			case controls.TypeBoolean:
				var p wot_properties.BooleanPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			case controls.TypeInteger:
				var p wot_properties.IntegerPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			case controls.TypeNumber:
				var p wot_properties.NumberPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			case controls.TypeObject:
				var p wot_properties.ObjectPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			case controls.TypeArray:
				var p wot_properties.ArrayPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			case controls.TypeNull:
				var p wot_properties.NullPropertyAffordance
				prop.ToVal(&p)
				if &p != nil {
					props[name] = &p
				}
			default:
				continue
			}
		}
		if props != nil && len(props) > 0 {
			t.Properties = props
		}
	}

	var actions ThingActions
	if a := json.Get(data, "actions"); a.LastError() == nil {
		a.ToVal(&actions)
		if actions != nil {
			t.Actions = actions
		}
	}

	var events ThingEvents
	if e := json.Get(data, "events"); e.LastError() == nil {
		e.ToVal(&events)
		if e != nil {
			t.Events = events
		}
	}

	var links []controls.Link
	if l := json.Get(data, "links"); l.LastError() == nil {
		l.ToVal(&links)
		if l != nil {
			t.Links = links
		}
	}

	var forms []controls.Form
	if cf := json.Get(data, "forms"); cf.LastError() == nil {
		cf.ToVal(&forms)
		if forms != nil {
			t.Forms = forms
		}
	}

	var security controls.ArrayOrString
	if ca := json.Get(data, "security"); ca.LastError() == nil {
		ca.ToVal(&security)
		if security != "" {
			t.Security = security
		}
	}

	var securityDefinitions ThingSecurityDefinitions
	if s := json.Get(data, "securityDefinitions"); s.LastError() == nil {
		s.ToVal(&securityDefinitions)
		if securityDefinitions != nil {
			t.SecurityDefinitions = securityDefinitions
		}
	}

	var profile []controls.URI
	if pf := json.Get(data, "profile"); pf.LastError() == nil {
		pf.ToVal(&profile)
		if profile != nil {
			t.Profile = profile
		}
	}

	var schemaDefinitions ThingSchemaDefinitions
	if tsd := json.Get(data, "schemaDefinitions"); tsd.LastError() == nil {
		tsd.ToVal(&schemaDefinitions)
		if securityDefinitions != nil {
			t.SchemaDefinitions = schemaDefinitions
		}
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
