package models

import (
	"addon"
	"addon/wot"
	"fmt"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"strings"
)

type Thing struct {
	AtContext   []string `json:"@context"`
	Title       string   `json:"title"`
	ID          string   `json:"id"`
	AtType      []string `json:"@type"`
	Description string   `json:"description,omitempty"`

	Properties map[string]*Property `json:"properties"`
	Actions    map[string]*Action   `json:"actions,omitempty"`
	Events     map[string]*Event    `json:"events,omitempty"`
	Forms      []wot.Form           `json:"forms,omitempty"`

	//The configuration  of the device
	Pin                 addon.PIN
	CredentialsRequired bool `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconData           string `json:"iconData,omitempty"`
}

func NewThing(description string) (thing *Thing) {

	id := gjson.Get(description, "id").String()
	if id == "" {
		return nil
	}
	title := gjson.Get(description, "title").String()
	if title == "" {
		title = id
	}
	thingId := fmt.Sprintf("/things/%s", id)

	var atContext []string
	for _, c := range gjson.Get(description, "@context").Array() {
		atContext = append(atContext, c.String())
	}

	var atType []string
	for _, c := range gjson.Get(description, "@type").Array() {
		atType = append(atType, c.String())
	}
	if len(atType) < 1 {
		log.Info("new thing err: @type")
		return nil
	}

	t := &Thing{
		AtContext: atContext,
		AtType:    atType,
		Title:     title,
		ID:        thingId,
	}

	t.IconData = gjson.Get(description, "iconData").String()
	t.Connected = gjson.Get(description, "connected").Bool()

	if gjson.Get(description, "description").Exists() {
		t.Description = gjson.Get(description, "description").String()
	}

	if gjson.Get(description, "credentialsRequired").Exists() {
		t.CredentialsRequired = gjson.Get(description, "credentialsRequired").Bool()
	}

	if gjson.Get(description, "selectedCapability").Exists() {
		sc := gjson.Get(description, "selectedCapability").String()
		for _, s := range t.AtType {
			if s == sc {
				t.SelectedCapability = sc
			}
		}
	} else {
		if len(t.AtType) > 0 {
			t.SelectedCapability = t.AtType[0]
		}
	}

	if gjson.Get(description, "pin").Exists() {
		var pin addon.PIN
		pin.Required = gjson.Get(description, "pin.required").Bool()
		pin.Pattern = gjson.Get(description, "pin.pattern").Value()
		t.Pin = pin
	}
	var props = gjson.Get(description, "properties").Map()
	if len(props) > 0 {
		for name, data := range props {
			prop := NewProperty(data.String())
			if prop != nil {
				if prop.Forms == nil {
					prop.Forms = append(prop.Forms, wot.Form{
						Href: fmt.Sprintf("%s/properties/%s", t.ID, name),
					})
				}
				t.Properties[name] = prop
			}
		}
	}

	var actions = gjson.Get(description, "actions").Map()
	if len(actions) > 0 {
		t.Actions = make(map[string]*Action)
		for name, a := range actions {
			var action Action
			err := json.UnmarshalFromString(a.String(), &action)
			if err != nil {
				continue
			}
			action.ID = a.Get("id").String()
			if action.Forms == nil {
				action.Forms = append(action.Forms, wot.Form{Href: fmt.Sprintf("%s/actions/%s", t.ID, name)})
			}
			t.Actions[name] = &action
		}

		t.Forms = append(t.Forms, wot.Form{Rel: "actions", Href: thingId + util.ActionsPath})
	}

	var events = gjson.Get(description, "events").Map()
	if len(events) > 0 {
		t.Events = make(map[string]*Event)
		for name, e := range events {
			var event Event
			err := json.UnmarshalFromString(e.String(), &event)
			if err != nil {
				continue
			}
			if !e.Get("id").Exists() {
				continue
			}
			event.ID = e.Get("id").String()
			if event.Forms == nil {
				event.Forms = append(event.Forms, wot.Form{Href: fmt.Sprintf("%s/events/%s", thingId, name)})
			}
			t.Events[name] = &event
		}
		t.Forms = append(t.Forms, wot.Form{Rel: "actions", Href: t.ID + util.EventsPath})
	}
	if t.Forms == nil {
		t.Forms = append(t.Forms, wot.Form{Rel: "alternate", ContentType: "text/html", Href: t.ID})
		t.Forms = append(t.Forms, wot.Form{Rel: "alternate", Href: fmt.Sprintf("wss://localhost/%s", thingId)})
	}

	return t
}

func (t *Thing) setSelectedCapability(sel string) {
	t.SelectedCapability = sel
	_ = t.update
	t.Publish(util.MODIFIED, t)
}

func (t *Thing) findProperty(propName string) (interface{}, error) {
	prop, ok := t.Properties[propName]
	if !ok {
		return nil, fmt.Errorf("thing(%s) can not found properties(%s)", t.ID, propName)
	}
	return prop, nil
}

func (t *Thing) GetProperty(propName string) interface{} {
	prop, ok := t.Properties[propName]
	if !ok {
		log.Debug("thing(%s) can not found properties(%s)", t.ID, propName)
		return nil
	}
	return prop
}

func (t *Thing) GetID() string {
	sl := strings.Split(t.ID, "/")
	id := sl[len(sl)-1]
	return id
}

func (t *Thing) SetTitle(title string) string {
	t.Title = title
	_ = t.update
	t.Publish(util.MODIFIED, t)
	return t.GetDescription()
}

func (t *Thing) GetTitle() string {
	return t.Title
}

func (t *Thing) AddAction(action *Action) error {
	return nil
}

func (t *Thing) SetSelectedCapability(selectedCapability string) {
	t.SelectedCapability = selectedCapability
}

func (t *Thing) setConnected(connected bool) {

	t.Publish(util.CONNECTED, connected)
	t.Connected = connected

}

func (t *Thing) IsConnected() bool {
	return t.Connected
}

func (t *Thing) RemoveAction(a *Action) bool {
	_, ok := t.Actions[a.Name]
	return ok
}

func (t *Thing) GetDescription() string {
	s, err := json.MarshalToString(t)
	if err != nil {
		return ""
	}
	return s
}

func (t *Thing) save() (err error) {
	return database.SetSetting(t.ID, t.GetDescription())
}

func (t *Thing) update(thing *Thing) {
	t.AtContext = thing.AtContext
	t.AtType = thing.AtType
	t.Description = thing.Description
	t.Properties = thing.Properties
	t.Actions = thing.Actions
	t.Events = thing.Events

	if thing.SelectedCapability != "" {
		for _, s := range t.AtType {
			if s == thing.SelectedCapability {
				t.setSelectedCapability(thing.SelectedCapability)
			}
			break
		}
	}
	_ = database.UpdateThing(t.GetID(), t.GetDescription())
}

func (t *Thing) Subscribe(typ string, f interface{}) {
	go func() {
		err := bus.Subscribe(t.GetID()+"."+typ, f)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (t *Thing) Unsubscribe(typ string, f interface{}) {
	go func() {
		err := bus.Subscribe(t.GetID()+"."+typ, f)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (t *Thing) Publish(typ string, args ...interface{}) {
	go bus.Publish(t.GetID()+"."+typ, args...)
}
