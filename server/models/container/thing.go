package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/bus/topic"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type ThingPin struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern"`
}

type eventBus interface {
	Pub(topic.Topic, ...any)
}

type Thing struct {
	*wot.Thing
	Pin                 *ThingPin `json:"pin,omitempty"`
	CredentialsRequired bool      `json:"credentialsRequired,omitempty"`

	//The state  of the thing
	SelectedCapability string `json:"selectedCapability,omitempty"`
	Connected          bool   `json:"-"`

	//FloorplanVisibility bool `json:"floorplanVisibility"`
	//FloorplanX          uint `json:"floorplanX"`
	//FloorplanY          uint `json:"floorplanY"`
	//LayoutIndex         uint `json:"layoutIndex"`
	GroupId string `json:"groupId,omitempty"`

	container *ThingsContainer
	bus       eventBus
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	var thing wot.Thing
	err := json.Unmarshal(data, &thing)
	if err != nil {
		return err
	}
	t.Thing = &thing

	var pin ThingPin
	if p := json.Get(data, "pin"); p.LastError() == nil {
		p.ToVal(&pin)
		if &pin != nil {
			t.Pin = &pin
		}
	}

	t.CredentialsRequired = json.Get(data, "credentialsRequired").ToBool()
	t.Connected = json.Get(data, "connected").ToBool()
	t.GroupId = json.Get(data, "groupId").ToString()
	t.SelectedCapability = json.Get(data, "selectedCapability").ToString()
	if t.SelectedCapability == "" {
		t.SelectedCapability = t.AtType[0]
	}
	var b = false

	for _, s := range t.AtType {
		if s == t.SelectedCapability {
			b = true
		}
	}
	if !b {
		return fmt.Errorf("selectedCapability err")
	}
	return nil
}

func (t *Thing) SetSelectedCapability(selectedCapability string) bool {
	for _, s := range t.AtType {
		if s == selectedCapability {
			t.SelectedCapability = selectedCapability
			t.bus.Pub(topic.ThingModify, t.GetId())
			return true
		}
	}
	return false
}

func (t *Thing) SetTitle(title string) bool {
	if t.Title == title {
		return false
	}
	t.Title = title
	t.bus.Pub(topic.ThingModify, t.GetId())
	return true
}

func (t *Thing) setConnected(connected bool) {
	if t.Connected == connected {
		return
	}
	t.Connected = connected
	t.bus.Pub(topic.ThingConnected, t.GetId(), connected)
}

func (t *Thing) remove() {
	t.bus.Pub(topic.ThingRemoved, t.GetId())
}

func (t *Thing) added() {
	t.bus.Pub(topic.ThingAdded, t.GetId())
}

func (t *Thing) AddAction(name string) bool {
	_, ok := t.Actions[name]
	return ok
}

func (t *Thing) RemoveAction(name string) bool {
	_, ok := t.Actions[name]
	return ok
}
