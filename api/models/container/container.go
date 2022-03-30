package container

import (
	"context"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/gofiber/fiber/v2"
	json "github.com/json-iterator/go"
	"sync"
	"time"
)

type Container interface {
	SetThingPropertyValue(thingId, propertyName string, value any) (any, error)
	GetThingPropertyValue(thingId, propertyName string) (any, error)
	Publish(topic2 topic.Topic, args ...any)
	Subscribe(topic2 topic.Topic, f any) func()
	GetThing(thing string) *Thing
	GetThings() []*Thing
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
	DeleteThing(id string) error
	CreateThing(id string, thing *Thing) error
	UpdateThing(id string, thing *Thing) error
	GetThings() map[string][]byte
}

type Listener interface {
	func(message topic.ThingAddedMessage)
	func(message topic.ThingRemovedMessage)
	func(message topic.ThingAddedMessage)
	func(message topic.ThingModifyMessage)
	func(message topic.ThingPropertyChangedMessage)
}

// ThingsContainer 管理所有Thing Models
// 向外部发送ThingAdded,ThingRemoved,ThingModify,ThingConnected事件
type ThingsContainer struct {
	things  sync.Map
	manager Manager
	store   ThingsStorage
	logger  logging.Logger
	bus.ThingsBus
	Listeners map[string]Listener
	//ThingAddedFuncs           map[string]func(message topic.ThingAddedMessage)
	//ThingRemovedFuncs         map[string]func(message topic.ThingRemovedMessage)
	//ThingConnectedFuncs       map[string]func(message topic.ThingAddedMessage)
	//ThingModifyFuncs          map[string]func(message topic.ThingModifyMessage)
	//ThingPropertyChangedFuncs map[string]func(message topic.ThingPropertyChangedMessage)
}

// NewThingsContainerModel 管理所有Thing Models
func NewThingsContainerModel(manager Manager, store ThingsStorage, log logging.Logger) *ThingsContainer {
	t := &ThingsContainer{}
	t.store = store
	t.manager = manager
	t.ThingsBus = bus.NewBus()
	t.logger = log
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
	a, _ := c.things.Load(id)
	t, ok := a.(*Thing)
	if !ok {
		return nil
	}
	return t
}

func (c *ThingsContainer) SetThingPropertyValue(thingId, propName string, value any) (any, error) {
	thing := c.GetThing(thingId)
	if thing == nil {
		return nil, util.NotFoundError("thing: %s not found", thing.GetId())
	}
	prop, ok := thing.Properties[propName]
	if !ok {
		return nil, util.NotFoundError("thing property:%s not found", propName)
	}
	if prop.IsReadOnly() {
		return nil, fiber.NewError(fiber.StatusBadRequest, "property read only")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancelFunc()
	}()
	v, err := c.manager.SetPropertyValue(ctx, thingId, propName, value)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *ThingsContainer) GetThings() []*Thing {
	return c.getThings()
}

func (c *ThingsContainer) getThings() (ts []*Thing) {
	ts = make([]*Thing, 0)
	c.things.Range(func(key, value any) bool {
		t, ok := value.(*Thing)
		if ok {
			ts = append(ts, t)
		}
		return true
	})
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

// CreateThing container和store中创建Thing
//ThingAdded事件
func (c *ThingsContainer) CreateThing(data []byte) (*Thing, error) {

	thing, err := NewThing(data, c)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err = c.store.CreateThing(thing.GetId(), thing)
	if err != nil {
		return nil, err
	}
	c.things.Store(thing.GetId(), thing)
	c.Publish(topic.ThingAdded, topic.ThingAddedMessage{
		ThingId: thing.GetId(),
		Data:    data,
	})
	return thing, nil
}

// RemoveThing 从container和store中删除Thing
//向定阅发送ThingRemoved事件
func (c *ThingsContainer) RemoveThing(thingId string) error {
	t := c.GetThing(thingId)
	err := c.store.DeleteThing(thingId)
	if err != nil {
		c.logger.Errorf(err.Error())
	}
	if t == nil {
		return util.NotFoundError("Thing Not Found")
	}
	c.things.Delete(thingId)
	c.Publish(topic.ThingRemoved, topic.ThingRemovedMessage{ThingId: thingId})
	return nil
}

func (c *ThingsContainer) UpdateThing(data []byte) error {
	thingId := json.Get(data, "id").ToString()
	if thingId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Thing id invalid")
	}
	thing := c.GetThing(thingId)
	if thing == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Thing Not Found")
	}
	newThing, err := NewThing(data, c)
	if err != nil {
		return err
	}
	return c.updateThing(newThing)
}

//在container和Store中更新Thing
func (c *ThingsContainer) updateThing(thing *Thing) error {
	err := c.store.UpdateThing(thing.GetId(), thing)
	if err != nil {
		return err
	}
	c.things.Store(thing.GetId(), thing)
	c.Publish(topic.ThingModify, topic.ThingModifyMessage{ThingId: thing.GetId()})
	return nil
}

//如果container为空，则从Store里加载所有的Things
func (c *ThingsContainer) updateThings() {
	things := c.getThings()
	if len(things) == 0 {
		for id, data := range c.store.GetThings() {
			thing, err := NewThing(data, c)
			if err != nil {
				c.logger.Errorf(err.Error())
				continue
			}
			c.things.Store(id, thing)
		}
	}
}

//当Manager中有设备移除，则更新Container中Thing中Connected为false
//并且向定阅发达ThingConnected通知
func (c *ThingsContainer) handleDeviceRemoved(thingId string) {
	t := c.GetThing(thingId)
	if t != nil {
		t.setConnected(false)
	}
	c.Publish(topic.ThingConnected, topic.ThingConnectedMessage{
		ThingId:   t.GetId(),
		Connected: true,
	})
}

//当Manager中有新设备，则更新Container中Thing状态
//并且向定阅发达ThingConnected通知
func (c *ThingsContainer) handleDeviceAdded(msg topic.DeviceAddedMessage) {
	t := c.GetThing(msg.DeviceId)
	if t != nil {
		t.setConnected(true)
	}
	c.Publish(topic.ThingConnected, topic.ThingConnectedMessage{
		ThingId:   msg.Id,
		Connected: true,
	})
}

//当Manager中有设备离线消息时，同步Thing的状态
//并向Container中的定阅发送ThingConnected消息
func (c *ThingsContainer) handleDeviceConnected(msg topic.DeviceConnectedMessage) {
	t := c.GetThing(msg.DeviceId)
	if t != nil {
		t.setConnected(msg.Connected)
	}
	c.Publish(topic.ThingConnected, topic.ThingConnectedMessage{
		ThingId:   t.GetId(),
		Connected: msg.Connected,
	})
}

func (c *ThingsContainer) handleDevicePropertyChanged(message topic.DevicePropertyChangedMessage) {
	t := c.GetThing(message.DeviceId)
	if t == nil {
		return
	}
	if p := t.GetProperty(message.PropertyName); p == nil {
		return
	}
	c.Publish(topic.ThingPropertyChanged, topic.ThingPropertyChangedMessage{
		ThingId:      t.GetId(),
		PropertyName: message.Name,
		Value:        message.Value,
	})
}

func (c *ThingsContainer) handleDeviceActionStatus(msg topic.DeviceActionStatusMessage) {
	t := c.GetThing(msg.DeviceId)
	if t == nil {
		return
	}
	t.onActionStatus(msg.Action)
}

func (c *ThingsContainer) handleDeviceEvent(deviceId string, event *events.EventDescription) {
	t := c.GetThing(deviceId)
	if t == nil {
		return
	}
	t.OnEvent(event)
}
