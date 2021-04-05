package models

import (
	"addon"
	"addon/wot"
	"fmt"
	"gateway/pkg/bus"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
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
	Forms      []map[string]string  `json:"forms,omitempty"`

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
		t.Properties = make(map[string]*Property)
		for name, p := range props {
			var property Property
			err := json.UnmarshalFromString(p.String(), &property)
			if err != nil {
				continue
			}
			property.Forms = append(property.Forms, wot.NewForm("href", fmt.Sprintf("%s/properties/%s", t.ID, name)))
			t.Properties[name] = &property
		}

		t.Forms = append(t.Forms, wot.NewForm("rel", "properties", "href", t.ID+util.PropertiesPath))
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
			action.Forms = append(action.Forms, wot.NewForm("href", fmt.Sprintf("%s/actions/%s", t.ID, name)))
			t.Actions[name] = &action
		}

		t.Forms = append(t.Forms, wot.NewForm("rel", "actions", "href", t.ID+util.ActionsPath))
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

			event.Forms = append(event.Forms, wot.NewForm("href", fmt.Sprintf("%s/events/%s", t.ID, name)))
			t.Events[name] = &event
		}
		t.Forms = append(t.Forms, wot.NewForm("rel", "actions", "href", t.ID+util.EventsPath))
	}
	t.Forms = append(t.Forms, wot.NewForm("rel", "alternate", "mediaType", "text/html", "href", t.ID))
	t.Forms = append(t.Forms, wot.NewForm("rel", "alternate", "href", fmt.Sprintf("wss://localhost/%s", t.ID)))

	return t
}

//func UnmarshalDevice(d []byte) (*addon.Device, error) {
//
//	id := json.Get(d, "id").ToString()
//	if id == "" {
//		return nil, fmt.Errorf("device id lost")
//	}
//	title := json.Get(d, "title").ToString()
//	if title == "" {
//		title = id
//	}
//
//	atContext := json.Get(d, "@context").Keys()
//	if len(atContext) == 0 {
//		t := json.Get(d, "@context").ToString()
//		if t != "" {
//			atContext = append(atContext, t)
//		}
//	}
//
//	var atType []string
//	json.Get(d, `@type`).ToVal(&atType)
//	if len(atType) == 0 {
//		return nil, fmt.Errorf("@type lost")
//	}
//
//	var properties map[string]*addon.Property
//	json.Get(d, "properties").ToVal(&properties)
//
//	var actions map[string]*addon.Action
//	json.Get(d, "actions").ToVal(&actions)
//
//	var events map[string]*addon.Event
//	json.Get(d, "actions").ToVal(&events)
//
//	var pin *addon.PIN
//	json.Get(d, "pin").ToVal(&pin)
//
//	device := &addon.Device{
//		ID:                  id,
//		AtContext:           atContext,
//		Title:               title,
//		AtType:              atType,
//		Description:         json.Get(d, "description").ToString(),
//		CredentialsRequired: json.Get(d, "credentialsRequired").ToBool(),
//		Pin:                 addon.PIN{},
//		AdapterId:           json.Get(d, "adapterId").ToString(),
//	}
//	if len(properties) > 0 {
//		device.Properties = make(map[string]addon.IProperty)
//		for n, p := range properties {
//			p.DeviceId = id
//			device.Properties[n] = p
//		}
//	}
//
//	if len(events) > 0 {
//		device.Events = make(map[string]*addon.Event)
//		for n, e := range events {
//			e.DeviceId = id
//			device.Events[n] = e
//		}
//	}
//
//	if len(actions) > 0 {
//		device.Actions = make(map[string]*addon.Action)
//		for n, a := range actions {
//			a.DeviceId = id
//			device.Actions[n] = a
//		}
//	}
//
//	if pin != nil {
//		device.Pin = *pin
//	}
//
//	return device, nil
//
//	id := gjson.Get(description, "id").String()
//	if id == "" {
//		return nil
//	}
//	th.ID = fmt.
//
//
//		err := json.Unmarshal(description, &th)
//	if len(th.AtContext) == 0 {
//		th.AtContext = []string{"https://webthings.io/schemas/"}
//	}
//	if err != nil {
//		return nil
//	}
//	return &th
//}

func (t *Thing) setSelectedCapability(sel string) {
	t.SelectedCapability = sel
	_ = t.update
	t.Publish(util.MODIFIED, t)
}

func (t *Thing) findProperty(propName string) (*Property, error) {
	prop, ok := t.Properties[propName]
	if !ok {
		return nil, fmt.Errorf("thing(%s) can not found properties(%s)", t.ID, propName)
	}
	return prop, nil
}

func (t *Thing) GetProperty(propName string) *Property {
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

//thing save to database must do this:
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

	if t.SelectedCapability != "" {
		for _, s := range t.AtType {
			if s == t.SelectedCapability {
				break
			}
			break
		}
		t.SelectedCapability = ""
	}

	_ = database.UpdateThing(t.GetID(), t.GetDescription())
}

func (t *Thing) Subscribe(typ string, f interface{}) {
	go bus.Subscribe(t.GetID()+"."+typ, f)
}

func (t *Thing) Unsubscribe(typ string, f interface{}) {
	go bus.Subscribe(t.GetID()+"."+typ, f)
}

func (t *Thing) Publish(typ string, args ...interface{}) {
	go bus.Publish(t.GetID()+"."+typ, args...)
}
