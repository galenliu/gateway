package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
)

type ThingsContainer interface {
	Container
}

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
	GetMapThings() map[string]*Thing
	CreateThing(data []byte) (*Thing, error)
	RemoveThing(id string) error
	UpdateThing(data []byte) error
}

type ThingsModel struct {
	things map[string]*Thing
	store  ThingsStorage

	logger logging.Logger
}

func NewThingsContainerModel(store ThingsStorage, log logging.Logger) *ThingsModel {
	t := &ThingsModel{}

	t.store = store
	t.logger = log
	t.things = make(map[string]*Thing)
	return t
}

func (c *ThingsModel) GetThing(id string) *Thing {
	t, ok := c.things[id]
	if !ok {
		return nil
	}
	return t
}

func (c *ThingsModel) GetThings() (ts []*Thing) {
	c.updateThings()
	for _, t := range c.things {
		ts = append(ts, t)
	}
	return
}

func (c *ThingsModel) GetMapThings() map[string]*Thing {
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

func (c *ThingsModel) CreateThing(data []byte) (*Thing, error) {
	t, err := c.handleCreateThing(data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *ThingsModel) RemoveThing(thingId string) error {
	err := c.handleRemoveThing(thingId)
	if err != nil {
		return err
	}
	return nil
}

func (c *ThingsModel) UpdateThing(data []byte) error {
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

func (c *ThingsModel) handleCreateThing(data []byte) (*Thing, error) {
	thingId := json.Get(data, "id").ToString()
	th, err := NewThingFromString(thingId, string(data))
	if err != nil {
		return nil, err
	}
	_, ok := c.things[th.GetID()]
	if ok {
		return nil, fmt.Errorf("thing id: %s is exited", th.GetID())
	}

	err = c.store.CreateThing(th.GetID(), th)
	if err != nil {
		return nil, err
	}
	c.things[th.GetID()] = th
	return c.things[th.GetID()], nil
}

func (c *ThingsModel) handleRemoveThing(thingId string) error {
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

func (c *ThingsModel) handleUpdateThing(data []byte) error {
	thingId := json.Get(data, "id").ToString()
	if _, ok := c.things[thingId]; ok {
		newThing, err := NewThingFromString(thingId, string(data))
		if err != nil {
			return err
		}
		if newThing != nil {
			c.things[newThing.ID.GetID()] = newThing
			err := c.store.UpdateThing(newThing.GetID(), newThing)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ThingsModel) updateThings() {
	if len(c.things) < 1 {
		for id, bytes := range c.store.GetThings() {
			thing, err := NewThingFromString(id, string(bytes))
			if err != nil {
				continue
			}
			c.things[id] = thing
		}
	}
}
