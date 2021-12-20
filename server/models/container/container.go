package container

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	bus "github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"net/http"
	"time"
)

type ThingsManager interface {
	SetPropertyValue(ctx context.Context, thingId, propertyName string, value any) (any, error)
	GetPropertyValue(thingId, propertyName string) (any, error)
	GetPropertiesValue(thingId string) (map[string]any, error)
}

// ThingsStorage CRUD
type ThingsStorage interface {
	RemoveThing(id string) error
	CreateThing(id string, thing any) error
	UpdateThing(id string, thing any) error
	GetThings() map[string]*Thing
}

type ThingsContainer struct {
	things  map[string]*Thing
	manager ThingsManager
	store   ThingsStorage
	logger  logging.Logger
	bus     containerBus
}

type containerBus interface {
	Sub(topic topic.Topic, fn any) func()
	Pub(topic topic.Topic, args ...any)
}

func NewThingsContainerModel(manager ThingsManager, store ThingsStorage, b *bus.Bus, log logging.Logger) *ThingsContainer {
	t := &ThingsContainer{}
	t.store = store
	t.manager = manager
	t.logger = log
	t.bus = b
	t.things = make(map[string]*Thing, 20)
	t.updateThings()
	_ = b.Sub(topic.DeviceAdded, t.handleDeviceAdded)
	_ = b.Sub(topic.DeviceRemoved, t.handleDeviceRemoved)
	_ = b.Sub(topic.DeviceConnected, t.handleDeviceConnected)
	_ = b.Sub(topic.DevicePropertyChanged, t.handleDevicePropertyChanged)
	_ = b.Sub(topic.DeviceActionStatus, t.handleDeviceActionStatus)
	_ = b.Sub(topic.DeviceEvent, t.handleDeviceEvent)
	return t
}

func (c *ThingsContainer) GetThing(id string) *Thing {
	t, ok := c.things[id]
	if !ok {
		return nil
	}
	return t
}

func (c *ThingsContainer) SetThingProperty(thingId, propName string, value any) (any, error) {
	thing := c.GetThing(thingId)
	if thing == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Thing not found")
	}
	prop, ok := thing.Properties[propName]
	if !ok {
		return nil, fiber.NewError(fiber.StatusNotFound, "property not found")
	}
	if prop.IsReadOnly() {
		return nil, fiber.NewError(fiber.StatusNotFound, "property read only")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancelFunc()
	}()
	v, err := c.manager.SetPropertyValue(ctx, thingId, propName, value)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return v, nil
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

	thingId := json.Get(data, "id").ToString()
	var thing Thing
	err := json.Unmarshal(data, &thing)
	if err != nil || thingId == "" {
		return nil, err
	}
	t := c.GetThing(thingId)
	if t != nil {
		return nil, fmt.Errorf("thing: %s is exited", t.GetId())
	}
	thing.container = c
	thing.bus = c.bus
	err = c.store.CreateThing(thingId, thing)
	if err != nil {
		return nil, err
	}
	id := thing.GetId()
	th := &thing
	c.things[id] = th
	return th, nil
}

func (c *ThingsContainer) RemoveThing(thingId string) error {
	err := c.store.RemoveThing(thingId)
	if err != nil {
		c.logger.Error("remove thing id: %s from Store err: %s", thingId, err.Error())
	}
	t, ok := c.things[thingId]
	if !ok || t == nil {
		return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("container has not thing id: %s", thingId))
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

func (c *ThingsContainer) handleUpdateThing(data []byte) error {
	thingId := json.Get(data, "id").ToString()
	if _, ok := c.things[thingId]; ok {
		var newThing Thing
		err := json.Unmarshal(data, &newThing)
		if err != nil {
			return err
		}
		c.things[newThing.Id.GetId()] = &newThing
		err = c.store.UpdateThing(newThing.GetId(), newThing)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *ThingsContainer) updateThings() {
	if len(c.things) < 1 {
		for id, thing := range c.store.GetThings() {
			thing.container = c
			thing.bus = c.bus
			c.things[id] = thing
		}
	}
}

func (c *ThingsContainer) handleDeviceRemoved(thingId string) {
	t := c.GetThing(thingId)
	if t != nil {
		t.setConnected(false)
	}
}

func (c *ThingsContainer) handleDeviceAdded(deviceId string, _ *addon.Device) {
	t := c.GetThing(deviceId)
	if t != nil {
		t.setConnected(true)
	}
}

func (c *ThingsContainer) handleDeviceConnected(deviceId string, connected bool) {
	t := c.GetThing(deviceId)
	if t != nil {
		t.setConnected(connected)
	}
}

func (c *ThingsContainer) handleDevicePropertyChanged(deviceId string, property *addon.PropertyDescription) {
	t := c.GetThing(deviceId)
	if t != nil {
		t.bus.Pub(topic.ThingPropertyChanged, t.GetId(), property)
	}
}

func (c *ThingsContainer) handleDeviceActionStatus(deviceId string, action *addon.ActionDescription) {
	t := c.GetThing(deviceId)
	if t != nil {
		t.bus.Pub(topic.ThingActionStatus, t.GetId(), action)
	}
}

func (c *ThingsContainer) handleDeviceEvent(deviceId string, event *addon.Event) {
	t := c.GetThing(deviceId)
	if t != nil {
		t.bus.Pub(topic.ThingEvent, t.GetId(), event)
	}
}
