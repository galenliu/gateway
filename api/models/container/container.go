package container

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/logging"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"time"
)

type Container interface {
	SetThingPropertyValue(thingId, propertyName string, value any) (any, error)
	GetThingPropertyValue(thingId, propertyName string) (any, error)
	Publish(topic2 topic.Topic, args ...any)
	Subscribe(topic2 topic.Topic, f any) func()
	GetThing(thing string) *Thing
}

type Manager interface {
	SetPropertyValue(ctx context.Context, thingId, propertyName string, value any) (any, error)
	GetPropertyValue(thingId, propertyName string) (any, error)
	GetPropertiesValue(thingId string) (map[string]any, error)
	Publish(topic2 topic.Topic, args ...any)
	Subscribe(topic2 topic.Topic, f any) func()
}

// ThingsStorage CRUD
type ThingsStorage interface {
	RemoveThing(id string) error
	CreateThing(id string, thing Thing) error
	UpdateThing(id string, thing Thing) error
	GetThings() map[string]Thing
}

type ThingsContainer struct {
	things  map[string]*Thing
	manager Manager
	store   ThingsStorage
	logger  logging.Logger
	*bus.Controller
}

func NewThingsContainerModel(manager Manager, store ThingsStorage, log logging.Logger) *ThingsContainer {
	t := &ThingsContainer{}
	t.store = store
	t.manager = manager
	t.Controller = bus.NewBusController()
	t.logger = log
	t.things = make(map[string]*Thing, 20)
	t.updateThings()
	_ = manager.Subscribe(topic.DeviceAdded, t.handleDeviceAdded)
	_ = manager.Subscribe(topic.DeviceRemoved, t.handleDeviceRemoved)
	_ = manager.Subscribe(topic.DeviceConnected, t.handleDeviceConnected)
	_ = manager.Subscribe(topic.DevicePropertyChanged, t.handleDevicePropertyChanged)
	_ = manager.Subscribe(topic.DeviceActionStatus, t.handleDeviceActionStatus)
	_ = manager.Subscribe(topic.DeviceEvent, t.handleDeviceEvent)
	return t
}

func (c *ThingsContainer) GetThing(id string) *Thing {
	t, ok := c.things[id]
	if !ok {
		return nil
	}
	return t
}

func (c *ThingsContainer) SetThingPropertyValue(thingId, propName string, value any) (any, error) {
	thing := c.GetThing(thingId)
	if thing == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Thing not found")
	}
	prop, ok := thing.Properties[propName]
	if !ok {
		return nil, fiber.NewError(fiber.StatusNotFound, "property:%s not found", propName)
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

func (c *ThingsContainer) GetThingPropertyValue(thingId, propName string) (any, error) {
	return c.manager.GetPropertyValue(thingId, propName)
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

	var thing Thing
	err := json.Unmarshal(data, &thing)
	if err != nil {
		return nil, err
	}
	thing.container = c
	thing.Created = &controls.DataTime{Time: time.Now()}
	thing.Modified = &controls.DataTime{Time: time.Now()}
	th := &thing
	th.added()
	c.things[thing.GetId()] = th
	return th, nil
}

func (c *ThingsContainer) RemoveThing(thingId string) {

	t, _ := c.things[thingId]
	if t == nil {
		c.logger.Errorf("thing with id %s not found", thingId)
		return
	}
	t.removed()
	delete(c.things, thingId)
	return
}

func (c *ThingsContainer) UpdateThing(data []byte) error {
	id := json.Get(data, "id")
	if id.ValueType() != json.StringValue {
		return fmt.Errorf("thing id invaild")
	}
	thingId := json.Get(data, "id").ToString()
	if _, ok := c.things[thingId]; ok {
		var newThing Thing
		err := json.Unmarshal(data, &newThing)
		if err != nil {
			return err
		}
		newThing.container = c

		c.things[newThing.Id.GetId()] = &newThing
		newThing.update()
	}
	return nil
}

func (c *ThingsContainer) updateThings() {
	if len(c.things) < 1 {
		for id, thing := range c.store.GetThings() {
			thing.container = c
			c.things[id] = &thing
		}
	}
}

func (c *ThingsContainer) handleDeviceRemoved(thingId string) {
	t := c.GetThing(thingId)
	if t != nil {
		t.setConnected(false)
	}
}

func (c *ThingsContainer) handleDeviceAdded(msg topic.DeviceAddedMessage) {
	t := c.GetThing(msg.DeviceId)
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

func (c *ThingsContainer) handleDevicePropertyChanged(message topic.DevicePropertyChangedMessage) {
	t := c.GetThing(message.DeviceId)
	t.onPropertyChanged(message.PropertyDescription)
}

func (c *ThingsContainer) handleDeviceActionStatus(deviceId string, action *actions.ActionDescription) {
	t := c.GetThing(deviceId)
	t.onActionStatus(action)
}

func (c *ThingsContainer) handleDeviceEvent(deviceId string, event *events.EventDescription) {
	t := c.GetThing(deviceId)
	t.OnEvent(event)
}
