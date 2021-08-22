package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
)

// Container Things CRUD
type Container interface {
	GetThing(id string) *Thing
	GetThings() []*Thing
	GetMapThings() map[string]*Thing
	CreateThing(data []byte) (*Thing, error)
	RemoveThing(id string) error
	UpdateThing(data []byte) error
}

type ListenController interface {
	ListenCreateThing(func(data []byte) (*Thing, error))
	ListenRemoveThing(func(id string) error)
}

type FireController interface {
	FireThingAdded(thing *Thing)
	FireThingRemoved(id string)
}

type ThingsStorage interface {
	RemoveThing(id string) error
	CreateThing(id string, thing interface{}) error
	UpdateThing(id string, thing interface{}) error
	GetThings() map[string][]byte
}

type container struct {
	things map[string]*Thing
	store  ThingsStorage
	logger logging.Logger
}

func NewThingsContainer(store ThingsStorage, log logging.Logger) Container {
	instance := &container{}
	instance.store = store
	instance.logger = log
	instance.things = make(map[string]*Thing)
	return instance
}

func (c *container) GetThing(id string) *Thing {
	t, ok := c.things[id]
	if !ok {
		return nil
	}
	return t
}

func (c *container) GetThings() (ts []*Thing) {
	c.updateThings()
	for _, t := range c.things {
		ts = append(ts, t)
	}
	return
}

func (c *container) GetMapThings() map[string]*Thing {
	things := c.GetThings()
	if things == nil {
		return nil
	}
	var thingsMap = make(map[string]*Thing)
	for _, th := range things {
		thingsMap[th.GetID()] = th
	}
	return thingsMap
}

func (c *container) CreateThing(data []byte) (*Thing, error) {
	t, err := c.handleCreateThing(data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *container) RemoveThing(thingId string) error {
	err := c.handleRemoveThing(thingId)
	if err != nil {
		return err
	}
	return nil
}

func (c *container) UpdateThing(data []byte) error {
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

func (c *container) handleCreateThing(data []byte) (*Thing, error) {
	th, err := NewThingFromString(string(data))
	if err != nil {
		return nil, err
	}
	_, ok := c.things[th.GetID()]
	if ok {
		return nil, fmt.Errorf("thing id: %s is exited", th.GetID())
	}
	bytes, err := json.Marshal(th)
	err = c.store.CreateThing(th.GetID(), bytes)
	if err != nil {
		return nil, err
	}
	c.things[th.GetID()] = th
	return c.things[th.GetID()], nil
}

func (c *container) handleRemoveThing(thingId string) error {
	err := c.store.RemoveThing(thingId)
	if err != nil {
		c.logger.Error("remove thing id: %s from Store err: %s", thingId, err.Error())
	}
	_, ok := c.things[thingId]
	if !ok {
		c.logger.Info("container has not thing id: %s", thingId)
	}
	delete(c.things, thingId)
	return nil
}

func (c *container) handleUpdateThing(data []byte) error {
	thingId := json.Get(data, "id").ToString()
	if _, ok := c.things[thingId]; ok {
		newThing, err := NewThingFromString(string(data))
		if err != nil {
			return err
		}
		if newThing != nil {
			c.things[newThing.ID.GetID()] = newThing
			bytes, _ := json.Marshal(newThing)
			err := c.store.UpdateThing(newThing.GetID(), bytes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *container) updateThings() {
	if len(c.things) < 1 {
		for id, bytes := range c.store.GetThings() {
			thing, err := NewThingFromString(string(bytes))
			if err != nil {
				continue
			}
			c.things[id] = thing
		}
	}
}
