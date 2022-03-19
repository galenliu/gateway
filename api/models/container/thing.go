package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/util"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

	json "github.com/json-iterator/go"
	//"github.com/goccy/go-json"
	"sync"
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
	sync.Mutex
}

func NewThing(data []byte, container *ThingsContainer) (*Thing, error) {
	var thing Thing
	err := json.Unmarshal(data, &thing)
	if err != nil {
		return nil, err
	}
	thing.Created = &controls.DataTime{Time: time.Now()}
	thing.container = container
	thing.create()
	container.things[thing.GetId()] = &thing
	return &thing, nil
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	var thing wot.Thing
	err := json.Unmarshal(data, &thing)
	if err != nil {
		return err
	}
	t.Mutex = sync.Mutex{}
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
			err := t.container.store.UpdateThing(t.GetId(), t)
			if err != nil {
				return false
			}
			go t.container.Publish(topic.ThingModify, topic.ThingModifyMessage{ThingId: t.GetId()})
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
	err := t.container.store.UpdateThing(t.GetId(), t)
	if err != nil {
		return false
	}
	go t.container.Publish(topic.ThingModify, topic.ThingModifyMessage{ThingId: t.GetId()})
	return true
}

func (t *Thing) setConnected(connected bool) {
	if t.Connected == connected {
		return
	}
	t.Connected = connected
	go t.container.Publish(topic.ThingConnected, topic.ThingConnectedMessage{
		ThingId:   t.GetId(),
		Connected: connected,
	})
}

func (t *Thing) removed() {
	err := t.container.store.DeleteThing(t.GetId())
	if err != nil {
		fmt.Printf("delete thing err:%s", err.Error())
	}
	go t.container.Publish(topic.ThingRemoved, topic.ThingRemovedMessage{ThingId: t.GetId()})
}

func (t *Thing) update() {
	err := t.container.store.UpdateThing(t.GetId(), t)
	if err != nil {
		return
	}
	go t.container.Publish(topic.ThingModify, topic.ThingModifyMessage{ThingId: t.GetId()})
}

func (t *Thing) create() {
	err := t.container.store.CreateThing(t.GetId(), t)
	if err != nil {
		t.container.logger.Errorf(err.Error())
	}
	t.container.Publish(topic.ThingAdded, topic.ThingAddedMessage{ThingId: t.GetId(), Data: []byte(util.JsonIndent(t))})
}

func (t *Thing) onPropertyChanged(msg topic.DevicePropertyChangedMessage) {
	if p := t.GetProperty(msg.PropertyName); p == nil {
		return
	}
	t.container.Publish(topic.ThingPropertyChanged, topic.ThingPropertyChangedMessage{
		ThingId:      t.GetId(),
		PropertyName: msg.Name,
		Value:        msg.Value,
	})
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
	go t.container.Publish(topic.ThingEvent, t.GetId(), event)
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
	go t.container.Subscribe(topic.ThingEvent+topic.Topic(t.GetId()), f)
}
