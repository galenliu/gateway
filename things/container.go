package things

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	models2 "github.com/galenliu/gateway/pkg/wot/models"
	AddonManager "github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/wot/models"
)

//ThingsEventBus Things和bus通讯的接口
type ThingsEventBus interface {
	ListenController
	FireController
}

type ListenController interface {
	ListenCreateThing(func(data []byte) error)
	ListenRemoveThing(func(id string))
}

type FireController interface {
	FireThingAdded(thing *Thing)
	FireThingRemoved(id string)
}

type ThingsStore interface {
	RemoveThing(id string) error
	SaveThing(t *Thing) error
	GetThings() []string
}

type ThingsContainer interface {
	GetThing(id string) *models2.Thing
	GetThings() []*models2.Thing
	CreateThing(data []byte) error
}

type Options struct {
}

type container struct {
	things  map[string]*Thing
	Actions *models2.Actions
	options Options
	store   ThingsStore
	logger  logging.Logger
	bus     ThingsEventBus
}

func NewThingsContainer(option Options, store ThingsStore, bus ThingsEventBus, log logging.Logger) ThingsContainer {

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

func (c *container) CreateThing(data []byte) error {
	return c.handleCreateThing(data)
}

func (c *container) handleCreateThing(data []byte) error {
	th, err := models.NewThingFromString(string(data))
	if err != nil {
		return err
	}
	_, ok := c.things[th.GetID()]
	if ok {
		return fmt.Errorf("thing id: %s is exited", th.GetID())
	}
	c.things[th.GetID()] = th
	err = c.store.SaveThing(th)
	if err != nil {
		return err
	}
	c.bus.FireThingAdded(th)
	return nil
}

func (c *container) handleRemoveThing(thingId string) {
	err := c.store.RemoveThing(thingId)
	if err != nil {
		c.logger.Error("remove thing id: %s from store err: %s", thingId, err.Error())
	}
	_, ok := c.things[thingId]
	if ok {
		c.logger.Error("container has not thing id: %s", thingId)
		return
	}
	delete(c.things, thingId)
	c.bus.FireThingRemoved(thingId)
}

func (c *container) SetThingProperty(thingId, propName string, value interface{}) (interface{}, error) {
	var th = c.GetThing(thingId)
	if th == nil {
		return nil, fmt.Errorf("thing(%s) can not found", thingId)
	}
	return AddonManager.SetProperty(thingId, propName, value)

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
