package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/log"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/tidwall/gjson"

	"encoding/json"

	"time"
)

type ThingPin struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern"`
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
	GroupId   string `json:"groupId,omitempty"`
	container *ThingsContainer
}

func NewThing(data []byte, container *ThingsContainer) (*Thing, error) {
	var thing Thing
	err := json.Unmarshal(data, &thing)
	if err != nil {
		return nil, err
	}
	thing.Created = &controls.DataTime{Time: time.Now()}
	thing.container = container
	return &thing, nil
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	var thing wot.Thing
	err := json.Unmarshal(data, &thing)
	if err != nil {
		return err
	}
	if thing.Title == "" {
		return fmt.Errorf("thing title cannot be empty")
	}
	if thing.AtContext == "" {
		return fmt.Errorf("thing @context cannot be empty")
	}
	if thing.Id == "" {
		thing.Id = controls.NewURI(thing.Title)
	}
	t.Thing = &thing

	var pin ThingPin
	if result := gjson.GetBytes(data, "pin"); result.Exists() {
		err := json.Unmarshal(data[result.Index:result.Index+len(result.Raw)], &pin)
		if err != nil {
			t.Pin = &pin
		}
	}
	t.CredentialsRequired = gjson.GetBytes(data, "credentialsRequired").Bool()
	t.Connected = gjson.GetBytes(data, "connected").Bool()
	t.GroupId = gjson.GetBytes(data, "groupId").String()
	t.SelectedCapability = gjson.GetBytes(data, "selectedCapability").String()
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

func (t *Thing) SetSelectedCapability(selectedCapability string) {
	t.SelectedCapability = selectedCapability
	err := t.container.store.UpdateThing(t.GetId(), t)
	if err != nil {
		log.Errorf(err.Error())
	}

}

func (t *Thing) SetTitle(title string) {
	t.Title = title
	err := t.container.store.UpdateThing(t.GetId(), t)
	if err != nil {
		log.Errorf(err.Error())
	}
}

func (t *Thing) setConnected(connected bool) {
	t.Connected = connected
}

func (t *Thing) onActionStatus(a actions.ActionDescription) {
	go t.container.Publish(topic.ThingActionStatus, topic.ThingActionStatusMessage{
		ThingId: t.GetId(),
		Action: topic.ThingActionDescription{
			Id:            a.Id,
			Name:          a.Name,
			Input:         a.Input,
			Status:        a.Status,
			TimeRequested: a.TimeRequested,
			TimeCompleted: a.TimeCompleted,
		},
	})
}

func (t *Thing) OnEvent(event *events.EventDescription) {
	t.container.Publish(topic.ThingEvent, t.GetId(), event)
}

func (t *Thing) AddAction(name string) bool {
	_, ok := t.Actions[name]
	return ok
}

func (t *Thing) RemoveAction(name string) bool {
	_, ok := t.Actions[name]
	return ok
}

func (t *Thing) GetPropertyValue(name string) (any, error) {
	return t.container.manager.GetPropertyValue(t.GetId(), name)
}

func (t *Thing) AddEventSubscription(f func(message topic.ThingEventMessage)) {
	_ = t.container.Subscribe(topic.ThingEvent+topic.Topic(t.GetId()), f)
}
