package things

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
)

type eventBus interface {
	Publish(string, ...interface{})
	Subscribe(topic string, fn interface{})
}

// Container Things CRUD
type Container interface {
	GetThing(id string) *Thing
	GetThings() []*Thing
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
	bus     eventBus
}

func NewThingsContainer(option Options, store Store, bus eventBus, log logging.Logger) Container {

	instance := &container{}
	instance.options = option
	instance.store = store
	instance.bus = bus
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

func (c *container) CreateThing(data []byte) (*Thing, error) {
	t, err := c.handleCreateThing(data)
	if err != nil {
		return nil, err
	}
	c.bus.Publish(util.ThingCreated, t)
	return t, nil
}

func (c *container) RemoveThing(thingId string) error {
	err := c.handleRemoveThing(thingId)
	if err != nil {
		return err
	}
	c.bus.Publish(util.ThingRemoved, thingId)
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
	c.bus.Publish(util.ThingModified, id.ToString())
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
