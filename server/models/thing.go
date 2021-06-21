package models

import (
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	data_schema2 "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	"github.com/galenliu/gateway/server/models/model"
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

	Properties map[string]*model.Property `json:"properties,omitempty"`
	Actions    map[string]*model.Action   `json:"_actions,omitempty"`
	Events     map[string]*model.Event    `json:"events,omitempty"`

	Forms   []hypermedia_controls.Form `json:"forms,omitempty"`
	Links   []hypermedia_controls.Link `json:"links,omitempty"`
	Support interface{}                `json:"support,omitempty"`

	Version  interface{} `json:"version,omitempty"`
	Created  *time.Time  `json:"created,omitempty"`
	Modified *time.Time  `json:"modified,omitempty"`

	Security interface{} `json:"security,omitempty"`
	//The configuration  of the device
	Pin                 *addon.PIN `json:"pin,omitempty"`
	CredentialsRequired bool       `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability"`
	Connected          bool   `json:"connected"`
	IconData           string `json:"iconData,omitempty"`
}

// NewThingFromString 把传入description组装成一个thing对象
func NewThingFromString(description string) (thing *Thing) {

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
		logging.Info("new thing err: @type")
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
		t.Pin = &pin
	}

	if gjson.Get(description, "properties").Exists() {
		var props = gjson.Get(description, "properties").Map()
		if len(props) > 0 {
			t.Properties = make(map[string]*model.Property)
			for name, data := range props {
				prop := model.NewPropertyFromString(data.String())
				if prop != nil {
					if prop.Forms == nil {
						prop.Forms = append(prop.Forms, hypermedia_controls.Form{
							Href:        fmt.Sprintf("%s%s/%s", thingId, util.PropertiesPath, name),
							ContentType: data_schema2.ApplicationJson,
							Op:          []string{hypermedia_controls.ReadProperty, hypermedia_controls.WriteProperty},
						})
					}
					t.Properties[name] = prop
				}
			}
		}
		t.Forms = append(t.Forms, hypermedia_controls.Form{Op: []string{hypermedia_controls.ReadallProperties, hypermedia_controls.WriteAllProperties}, Href: thingId + util.PropertiesPath, ContentType: data_schema2.ApplicationJson})
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
					action.Forms = append(action.Forms, hypermedia_controls.Form{Href: fmt.Sprintf("%s/_actions/%s", t.ID, name)})
				}
				if action.Name == "" {
					action.Name = name
				}
				t.Actions[name] = &action
			}

			t.Forms = append(t.Forms, hypermedia_controls.Form{Rel: "_actions", Href: thingId + util.ActionsPath})
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
					event.Forms = append(event.Forms, hypermedia_controls.Form{Href: fmt.Sprintf("%s/events/%s", thingId, name)})
				}
				t.Events[name] = &event
			}
			t.Forms = append(t.Forms, hypermedia_controls.Form{Rel: "_actions", Href: t.ID + util.EventsPath})
		}
	}

	if t.Forms == nil {
		t.Forms = append(t.Forms, hypermedia_controls.Form{Rel: "alternate", ContentType: "text/html", Href: t.ID})
		t.Forms = append(t.Forms, hypermedia_controls.Form{Rel: "alternate", Href: fmt.Sprintf("wss://localhost/%s", thingId)})
	}
	return t
}

func (t *Thing) setSelectedCapability(s string) {
	if t.SelectedCapability == s {
		return
	}
	for _, typ := range t.AtType {
		if s == typ {
			t.SelectedCapability = s
			err := t.save()
			if err != nil {
				return
			}

		}
	}

}

func (t *Thing) GetSelectedCapability() string {
	return t.SelectedCapability
}

func (t *Thing) GetID() string {
	sl := strings.Split(t.ID, "/")
	id := sl[len(sl)-1]
	return id
}

func (t *Thing) SetTitle(title string) string {
	if t.Title != title {
		t.Title = title
		err := t.save()
		if err != nil {
			logging.Info(err.Error())
		}

	}
	return t.GetDescription()
}

func (t *Thing) GetTitle() string {
	return t.Title
}

func (t *Thing) setConnected(connected bool) {
	err := t.save()
	if err != nil {
		logging.Info(err.Error())
	}
	t.Connected = connected
}

func (t *Thing) isConnected() bool {
	return t.Connected
}

func (t *Thing) AddAction(a *Action) bool {
	_, ok := t.Actions[a.Name]
	return ok
}

func (t *Thing) RemoveAction(a *Action) error {
	_, ok := t.Actions[a.Name]
	if !ok {
		return fmt.Errorf("invalid action name :%s", a.Name)
	}
	return nil
}

func (t *Thing) GetDescription() string {
	s, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return ""
	}
	return string(s)
}

func (t *Thing) updateFromString(data string) {
	thing := NewThingFromString(data)
	t.SetTitle(thing.Title)
	t.Description = thing.Description
	if thing.SelectedCapability != "" {
		t.setSelectedCapability(thing.SelectedCapability)
	}
}

func (t *Thing) save() error {
	err := database.UpdateThing(t.GetID(), t.GetDescription())
	if err != nil {
		return err
	}
	return nil
}

//func (t *Thing) Subscribe(typ string, f interface{}) {
//	go func() {
//		err := event_bus.Subscribe(t.GetID()+"."+typ, f)
//		if err != nil {
//			logging.Error(err.Error())
//		}
//	}()
//}
//
//func (t *Thing) Unsubscribe(typ string, f interface{}) {
//	go func() {
//		err := event_bus.Subscribe(t.GetID()+"."+typ, f)
//		if err != nil {
//			logging.Error(err.Error())
//		}
//	}()
//}
//
//func (t *Thing) Publish(typ string, args ...interface{}) {
//	go event_bus.Publish(t.GetID()+"."+typ, args...)
//}
