package core

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/util"
	dataSchema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	securityScheme "github.com/galenliu/gateway/pkg/wot/definitions/security_scheme"
	json "github.com/json-iterator/go"
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

	if title := json.Get(data, "title").ToString(); title == "" {
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

	m := make(map[string]securityScheme.SecurityScheme)
	json.Get(data, "securityDefinitions").ToVal(&m)
	t.SecurityDefinitions = m

	var propMap map[string]string
	if json.Get(data, "properties").ToVal(&propMap); len(propMap) > 0 {
		t.Properties = make(map[string]PropertyAffordance)
		for name, data := range propMap {
			prop := NewPropertyAffordanceFromString(data)
			t.Properties[name] = prop
		}
		t.Forms = append(t.Forms, controls.Form{Op: []string{controls.ReadallProperties, controls.WriteAllProperties}, Href: string(t.ID + util.PropertiesPath), ContentType: dataSchema.ApplicationJson})
	}

	var actionMap map[string]string
	if json.Get(data, "properties").ToVal(&actionMap); len(propMap) > 0 {

		t.Actions = make(map[string]ActionAffordance)
		for name, data := range actionMap {
			action := NewActionAffordanceFromString(data)
			t.Actions[name] = action
		}
		t.Forms = append(t.Forms, controls.Form{Href: string(t.ID + util.ActionsPath)})
	}

	var eventMap map[string]string
	if json.Get(data, "properties").ToVal(&eventMap); len(eventMap) > 0 {

		t.Events = make(map[string]EventAffordance)
		for name, data := range eventMap {
			event := NewEventAffordanceFromString(data)
			t.Events[name] = event
		}
		t.Forms = append(t.Forms, controls.Form{Href: string(t.ID + util.ActionsPath)})
	}

	if t.Forms == nil {
		t.Forms = append(t.Forms, controls.Form{ContentType: "text/html", Href: string(t.ID)})
		t.Forms = append(t.Forms, controls.Form{Href: fmt.Sprintf("wss://localhost/%s", t.ID)})
	}
	return t, nil
}
