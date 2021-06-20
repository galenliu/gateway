package models

import (
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/server/models/model"
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
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
			t.Publish(util.MODIFIED, t)
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
		t.Publish(util.MODIFIED, t)
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
	t.Publish(util.CONNECTED, connected)
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

func (t *Thing) Subscribe(typ string, f interface{}) {
	go func() {
		err := event_bus.Subscribe(t.GetID()+"."+typ, f)
		if err != nil {
			logging.Error(err.Error())
		}
	}()
}

func (t *Thing) Unsubscribe(typ string, f interface{}) {
	go func() {
		err := event_bus.Subscribe(t.GetID()+"."+typ, f)
		if err != nil {
			logging.Error(err.Error())
		}
	}()
}

func (t *Thing) Publish(typ string, args ...interface{}) {
	go event_bus.Publish(t.GetID()+"."+typ, args...)
}
