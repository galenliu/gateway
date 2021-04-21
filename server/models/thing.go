package models

import (
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway-addon/wot"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type Thing struct {
	AtContext    []string `json:"@context"`
	Title        string   `json:"title"`
	Titles       []string `json:"titles,omitempty"`
	ID           string   `json:"id"`
	AtType       []string `json:"@type"`
	Description  string   `json:"description,omitempty"`
	Descriptions []string `json:"descriptions,omitempty"`

	Properties map[string]*Property `json:"properties,omitempty"`
	Actions    map[string]*Action   `json:"actions,omitempty"`
	Events     map[string]*Event    `json:"events,omitempty"`

	Forms   []wot.Form  `json:"forms,omitempty"`
	Links   []wot.Link  `json:"links,omitempty"`
	Support interface{} `json:"support,omitempty"`

	Version  interface{} `json:"version,omitempty"`
	Created  *time.Time  `json:"created,omitempty"`
	Modified *time.Time  `json:"modified,omitempty"`

	Security interface{} `json:"security,omitempty"`
	//The configuration  of the device
	Pin                 addon.PIN `json:"pin,omitempty"`
	CredentialsRequired bool      `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconData           string `json:"iconData,omitempty"`
}

// NewThing 把传入description组装成一个thing对象
func NewThing(description string) (thing *Thing) {

	id := gjson.Get(description, "id").String()
	if id == "" {
		return nil
	}
	title := gjson.Get(description, "title").String()
	if title == "" {
		title = id
	}

	sl := strings.Split(id, "/")
	tid := sl[len(sl)-1]
	thingId := fmt.Sprintf("%s/%s", util.ThingsPath, tid)

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
		t.Properties = make(map[string]*Property)
		for name, data := range props {
			prop := NewProperty(data.String())
			if prop != nil {
				if prop.InteractionAffordance == nil {
					prop.InteractionAffordance = new(wot.InteractionAffordance)
				}
				if prop.Forms == nil {
					prop.Forms = append(prop.Forms, wot.Form{
						Href: fmt.Sprintf("%s/properties/%s", thingId, name),
					})
				}
				if prop.Name == "" {
					prop.Name = name
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
			if action.Name == "" {
				action.Name = name
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
	t.Publish(util.MODIFIED, t)
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

func (t *Thing) update(data string) {
	thing := NewThing(data)
	t.AtContext = thing.AtContext
	t.AtType = thing.AtType
	t.Description = thing.Description
	t.Properties = thing.Properties
	t.Actions = thing.Actions
	t.Events = thing.Events

	if thing.SelectedCapability != "" {
		for _, s := range t.AtType {
			if s == thing.SelectedCapability {
				t.SelectedCapability = thing.SelectedCapability
			}
			break
		}
	}
	err := t.save()
	if err != nil {
		return
	}
	t.Publish(util.MODIFIED, t)
}

func (t *Thing) save() error {
	err := database.UpdateThing(t.GetID(), t.GetDescription())
	if err != nil {
		return err
	}
	return nil
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
