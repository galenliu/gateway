package thing

import (
	"fmt"
	"gateway/log"
	"gateway/pkg/database"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
)

type Thing struct {
	ID          string   `json:"id"`
	AtContext   string   `json:"@context,omitempty"`
	AtType      []string `json:"@type"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`

	Links    []string `json:"links,omitempty"`
	BaseHref string   `json:"baseHref,omitempty"`
	Pin      struct {
		Required bool        `json:"required,omitempty"`
		Pattern  interface{} `json:"pattern,omitempty"`
	} `json:"pin,omitempty"`

	Href                string `json:"href"`
	CredentialsRequired bool   `json:"credentialsRequired,omitempty"`
	SelectedCapability  string `json:"selectedCapability"`

	Properties map[string]*Property `json:"properties,omitempty"`
	Actions    map[string]*Action   `json:"actions,omitempty"`
	Events     map[string]*Event    `json:"events,omitempty"`

	onConnectedFuncs       map[interface{}]func(bool)
	onRemovedFuncs         map[interface{}]func(thing *Thing)
	onModifiedFuncs        map[interface{}]func(thing *Thing)
	onPropertyChangedFuncs map[interface{}]func(property *Property)

	Connected bool `json:"connected"`
}

func NewThing(id string, description []byte) (thing *Thing) {
	var th Thing
	th.ID = id
	err := json.Unmarshal(description, &th)
	if err != nil {
		return nil
	}
	return &th
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

func (t *Thing) GetId() string {
	return t.ID
}

func (t *Thing) SetTitle(title string) {
	t.Title = title
}

func (t *Thing) GetTitle() string {
	return t.Title
}

func (t *Thing) AddAction(action *Action) error {
	return nil
}

func (t *Thing) AddConnectedSubscription(key interface{}, f func(bool)) func() {
	if t.onConnectedFuncs == nil {
		t.onConnectedFuncs = make(map[interface{}]func(bool))
	}
	t.onConnectedFuncs[key] = f
	var removeFunc = func() {
		delete(t.onConnectedFuncs, key)
	}
	return removeFunc
}

func (t *Thing) AddRemovedSubscription(key interface{}, f func(thing *Thing)) func() {
	if t.onRemovedFuncs == nil {
		t.onRemovedFuncs = make(map[interface{}]func(thing *Thing))
	}
	t.onRemovedFuncs[key] = f
	var removeFunc = func() {
		delete(t.onRemovedFuncs, key)
	}
	return removeFunc
}

func (t *Thing) AddPropertyChangedSubscription(key interface{}, f func(*Property)) func() {
	if t.onPropertyChangedFuncs == nil {
		t.onPropertyChangedFuncs = make(map[interface{}]func(property *Property))
	}
	t.onPropertyChangedFuncs[key] = f
	var removeFunc = func() {
		delete(t.onPropertyChangedFuncs, key)
	}
	return removeFunc
}

func (t *Thing) AddModifiedSubscription(conn *websocket.Conn, f func(thing *Thing)) func() {
	if t.onModifiedFuncs == nil {
		t.onModifiedFuncs = make(map[interface{}]func(thing *Thing))
	}
	t.onModifiedFuncs[conn] = f
	var removeFunc = func() {
		delete(t.onModifiedFuncs, conn)
	}
	return removeFunc
}

func (t *Thing) SetSelectedCapability(selectedCapability string) {
	t.SelectedCapability = selectedCapability
}

func (t *Thing) SetConnected(connected bool) {
	t.Connected = connected
	for _, f := range t.onConnectedFuncs {
		f(connected)
	}
}

func (t *Thing) IsConnected() bool {
	return t.Connected
}

func (t *Thing) RemoveAction(a *Action) bool {
	_, ok := t.Actions[a.Name]
	return ok
}

func (t *Thing) SetThingProperty(propertyName string, value interface{}) (interface{}, error) {
	prop, err := t.findProperty(propertyName)
	if err != nil {
		return value, err
	}
	newValue, setErr := prop.SetValue(value)
	if setErr != nil {
		for _, f := range t.onPropertyChangedFuncs {
			f(prop)
		}
		return newValue, setErr
	}
	return value, setErr
}

//thing save to database must do this:
func (t *Thing) GetDescription() string {
	s, err := json.MarshalToString(t)
	if err != nil {
		return ""
	}
	return s
}

func (t *Thing) Save() (err error) {
	return database.SetSetting(t.ID, t.GetDescription())
}

func (t *Thing) Update(thing *Thing) {
	_ = database.UpdateThing(t.ID, thing.GetDescription())
}
