package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
)

// ThingsStorage CRUD
type ThingsStorage interface {
	RemoveThing(id string) error
	CreateThing(id string, thing interface{}) error
	UpdateThing(id string, thing interface{}) error
	GetThings() map[string][]byte
}

// Container  Things
type Container interface {
	GetThing(id string) *Thing
	GetThings() []*Thing
	GetMapOfThings() map[string]*Thing
	CreateThing(data []byte) (*Thing, error)
	RemoveThing(id string) error
	UpdateThing(data []byte) error
}

type ThingsContainer struct {
	things map[string]*Thing
	store  ThingsStorage
	logger logging.Logger
	bus    containerBus
}

type containerBus interface {
	AddDeviceRemovedSubscription(fn func(deviceId string)) func()
	AddDeviceAddedSubscription(fn func(device *addon.Device)) func()

	PublishThingConnected(thingId string, connected bool)
	PublishThingRemoved(thingId string)
}

func NewThingsContainerModel(store ThingsStorage, bus containerBus, log logging.Logger) *ThingsContainer {
	t := &ThingsContainer{}
	t.store = store
	t.logger = log
	t.things = make(map[string]*Thing)
	_ = bus.AddDeviceRemovedSubscription(t.handleDeviceRemoved)
	_ = bus.AddDeviceAddedSubscription(t.handleDeviceAdded)
	return t
}

func (c *ThingsContainer) GetThing(id string) *Thing {
	t, ok := c.things[id]
	if !ok {
		return nil
	}
	return t
}

func (c *ThingsContainer) GetThings() (ts []*Thing) {
	c.updateThings()
	for _, t := range c.things {
		ts = append(ts, t)
	}
	return
}

func (c *ThingsContainer) GetMapOfThings() map[string]*Thing {
	things := c.GetThings()
	if things == nil {
		return nil
	}
	var thingsMap = make(map[string]*Thing)
	for _, th := range things {
		thingsMap[th.GetId()] = th
	}
	return thingsMap
}

func (c *ThingsContainer) CreateThing(data []byte) (*Thing, error) {
	t, err := c.handleCreateThing(data)
	if err != nil {
		return nil, err
	}
	t.container = c
	c.things[t.GetId()] = t
	return t, nil
}

func (c *ThingsContainer) RemoveThing(thingId string) error {
	t, err := c.handleRemoveThing(thingId)
	if err != nil {
		return err
	}
	t.remove()
	delete(c.things, thingId)
	return nil
}

func (c *ThingsContainer) UpdateThing(data []byte) error {
	id := json.Get(data, "id")
	if id.ValueType() != json.StringValue {
		return fmt.Errorf("thing id invaild")
	}
	err := c.handleUpdateThing(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *ThingsContainer) handleCreateThing(data []byte) (*Thing, error) {
	thingId := json.Get(data, "id").ToString()
	var thing Thing
	err := json.Unmarshal(data, &thing)
	if err != nil || thingId == "" {
		return nil, err
	}
	t := c.GetThing(thingId)
	if t != nil {
		return nil, fmt.Errorf("thing id: %s is exited", t.GetId())
	}
	err = c.store.CreateThing(thingId, thing)
	if err != nil {
		return nil, err
	}
	return &thing, nil
}

func (c *ThingsContainer) handleRemoveThing(thingId string) (*Thing, error) {
	err := c.store.RemoveThing(thingId)
	if err != nil {
		c.logger.Error("remove thing id: %s from Store err: %s", thingId, err.Error())
	}
	t, ok := c.things[thingId]
	if !ok {
		return nil, fmt.Errorf("container has not thing id: %s", thingId)
	}
	return t, nil
}

func (c *ThingsContainer) handleUpdateThing(data []byte) error {
	thingId := json.Get(data, "id").ToString()
	if _, ok := c.things[thingId]; ok {
		newThing, err := NewThingFromString(thingId, string(data))
		if err != nil {
			return err
		}
		if newThing != nil {
			c.things[newThing.Id.GetId()] = newThing
			err := c.store.UpdateThing(newThing.GetId(), newThing)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ThingsContainer) updateThings() {
	if len(c.things) < 1 {
		for id, bytes := range c.store.GetThings() {
			var thing Thing
			err := json.Unmarshal(bytes, &thing)
			if err != nil {
				continue
			}
			thing.container = c
			c.things[id] = &thing
		}
	}
}

func (c *ThingsContainer) handleDeviceRemoved(deviceId string) {
	t := c.GetThing(deviceId)
	if t != nil {
		t.setConnected(false)
	}
}

func (c *ThingsContainer) handleDeviceAdded(device *addon.Device) {
	t := c.GetThing(device.GetId())
	if t != nil {
		t.setConnected(true)
	}
}
