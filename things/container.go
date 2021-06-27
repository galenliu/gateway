package things

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
)

type Container interface {
	GetThing(id string) *Thing
	GetThings() []*Thing
	CreateThing(data []byte) (*Thing, error)
}

//ThingsEventBus Things和bus通讯的接口
type ThingsEventBus interface {
	ListenController
	FireController
}

type ListenController interface {
	ListenCreateThing(func(data []byte) (*Thing, error))
	ListenRemoveThing(func(id string) error)
}

type FireController interface {
	FireThingAdded(thing *Thing)
	FireThingRemoved(id string)
}

type Store interface {
	RemoveThing(id string) error
	SaveThing(t *Thing) error
	UpdateThing(t *Thing) error
	GetThings() []string
}

type Options struct {
}

type container struct {
	things  map[string]*Thing
	options Options
	store   Store
	logger  logging.Logger
	bus     ThingsEventBus
}

func NewThingsContainer(option Options, store Store, bus ThingsEventBus, log logging.Logger) Container {

	instance := &container{}
	instance.options = option
	instance.store = store
	instance.bus = bus
	instance.logger = log
	instance.things = make(map[string]*Thing)
	instance.bus.ListenCreateThing(instance.handleCreateThing)
	instance.bus.ListenRemoveThing(instance.handleRemoveThing)
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

func (c *container) CreateThing(data []byte) (*Thing, error) {
	t, err := c.handleCreateThing(data)
	if err != nil {
		return nil, err
	}
	c.bus.FireThingAdded(t)
	return t, nil
}

func (c *container) RemoveThing(thingId string) error {
	err := c.handleRemoveThing(thingId)
	if err != nil {
		return err
	}
	c.bus.FireThingRemoved(thingId)
	return nil
}

func (c *container) UpdateThing(data []byte) error {
	return c.handleUpdateThing(data)
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
	c.things[th.GetID()] = th
	err = c.store.SaveThing(th)
	if err != nil {
		return nil, err
	}
	return c.things[th.GetID()], nil
}

func (c *container) handleRemoveThing(thingId string) error {
	err := c.store.RemoveThing(thingId)
	if err != nil {
		c.logger.Error("remove thing id: %s from store err: %s", thingId, err.Error())
	}
	_, ok := c.things[thingId]
	if ok {
		return fmt.Errorf("container has not thing id: %s", thingId)
	}
	delete(c.things, thingId)
	return nil
}

func (c *container) updateThings() {
	if len(c.things) < 1 {
		for _, s := range c.store.GetThings() {
			t, err := NewThingFromString(s)
			if err != nil {
				c.logger.Errorf("new thing for database err: %s", err.Error())
				continue
			}
			c.things[t.GetID()] = t
		}
	}
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
			err := c.store.UpdateThing(newThing)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//func (ts *container) Subscribe(typ string, f interface{}) {
//	_ = event_bus.Subscribe("Things."+typ, f)
//}
//
//func (ts *container) Unsubscribe(typ string, f interface{}) {
//	_ = event_bus.Unsubscribe("Things."+typ, f)
//}
//
//func (ts *container) Publish(typ string, args ...interface{}) {
//	event_bus.Publish("Things."+typ, args...)
//}
