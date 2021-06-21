package core

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/util"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
	"github.com/galenliu/gateway/server/models/model"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type ThingInterface interface {
}

type Thing struct {
	AtContext    []string     `json:"@context"`
	Title        string       `json:"title"`
	Titles       []string     `json:"titles,omitempty"`
	ID           controls.URI `json:"id"`
	AtType       []string     `json:"@type"`
	Description  string       `json:"description,omitempty"`
	Descriptions []string     `json:"descriptions,omitempty"`

	Forms []controls.Form `json:"forms,omitempty"`
	Links []controls.Link `json:"links,omitempty"`

	Support interface{} `json:"support,omitempty"`
	Base    interface{} `json:"base"`

	Version  VersionInfo `json:"version,omitempty"`
	Created  *time.Time  `json:"created,omitempty"`
	Modified *time.Time  `json:"modified,omitempty"`

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
	t.ID = controls.URI(json.Get(data, "id").ToString())

	title := json.Get(data, "title").ToString()
	if title == "" {
		title = string(t.ID)
	}

	if c := json.Get(data, "@context"); c.ValueType() == json.StringValue {
		t.AtContext = []string{c.ToString()}
	} else {
		var l []string
		json.Get(data, "@context").ToVal(&l)
		t.AtContext = l
	}

	if len(t.AtContext) < 1 {
		return nil, fmt.Errorf("@context is mandatory")
	}

	if c := json.Get(data, "security"); c.ValueType() == json.StringValue {
		t.Security = []string{c.ToString()}
	} else {
		var l []string
		json.Get(data, "security").ToVal(&l)
		t.Security = l
	}

	t.Description = json.Get(data, "description").ToString()

	if gjson.Get(description, "properties").Exists() {
		var props = gjson.Get(description, "properties").Map()
		if len(props) > 0 {
			t.Properties = make(map[string]core2.PropertyAffordance)
			for name, data := range props {
				prop := core.NewPropertyFromString(data.String())
				if prop != nil {
					if prop.Forms == nil {
						prop.Forms = append(prop.Forms, hypermedia_controls2.Form{
							Href:        fmt.Sprintf("%s%s/%s", thingId, util.PropertiesPath, name),
							ContentType: dataSchema.ApplicationJson,
							Op:          []string{hypermedia_controls2.ReadProperty, hypermedia_controls2.WriteProperty},
						})
					}
					t.Properties[name] = prop
				}
			}
		}
		t.Forms = append(t.Forms, hypermedia_controls2.Form{Op: []string{hypermedia_controls2.ReadallProperties, hypermedia_controls2.WriteAllProperties}, Href: thingId + util.PropertiesPath, ContentType: dataSchema.ApplicationJson})
	}

	if gjson.Get(description, "_actions").Exists() {
		var actions = gjson.Get(description, "_actions").Map()
		if len(actions) > 0 {
			t.Actions = make(map[string]*model.Action)
			for name, a := range actions {
				var action model.Action
				err := json.UnmarshalFromString(a.String(), &action)
				if err != nil {
					continue
				}
				action.ID = a.Get("id").String()
				if action.Forms == nil {
					action.Forms = append(action.Forms, hypermedia_controls2.Form{Href: fmt.Sprintf("%s/_actions/%s", t.ID, name)})
				}
				if action.Name == "" {
					action.Name = name
				}
				t.Actions[name] = &action
			}

			t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "_actions", Href: thingId + util.ActionsPath})
		}
	}

	if gjson.Get(description, "events").Exists() {
		var events = gjson.Get(description, "events").Map()
		if len(events) > 0 {
			t.Events = make(map[string]*model.Event)
			for name, e := range events {
				var event model.Event
				err := json.UnmarshalFromString(e.String(), &event)
				if err != nil {
					continue
				}
				if !e.Get("id").Exists() {
					continue
				}
				event.ID = e.Get("id").String()
				if event.Forms == nil {
					event.Forms = append(event.Forms, hypermedia_controls2.Form{Href: fmt.Sprintf("%s/events/%s", thingId, name)})
				}
				t.Events[name] = &event
			}
			t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "_actions", Href: t.ID + util.EventsPath})
		}
	}

	if t.Forms == nil {
		t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "alternate", ContentType: "text/html", Href: t.ID})
		t.Forms = append(t.Forms, hypermedia_controls2.Form{Rel: "alternate", Href: fmt.Sprintf("wss://localhost/%s", thingId)})
	}
	return t, nil
}

func (t *Thing) GetID() string {
	sl := strings.Split(t.ID, "/")
	tid := sl[len(sl)-1]
	return tid
}

func (t *Thing) GetThingID() string {
	return t.ID
}
