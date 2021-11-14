package bus

import (
	bus "github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/logging"
)

type Bus struct {
	bus.Bus
	logger logging.Logger
}

func NewBusController(log logging.Logger) *Bus {
	b := &Bus{}
	b.logger = log
	b.Bus = bus.New()
	return b
}

func (b *Bus) PublishPairingTimeout() {
	b.Publish(constant.PairingTimeout)
}

func (b *Bus) PublishConnected(thingId string, connected bool) {
	b.Publish(thingId+"."+constant.Connected, connected)
}

func (b *Bus) PublishActionStatus(action *addon.Action) {
	b.Publish(constant.ActionStatus, action)
}

func (b *Bus) PublishPropertyChanged(thingId string, property *addon.Property) {
	b.Publish(thingId+"."+constant.PropertyChanged, property)
}

func (b *Bus) PublishEvent(event *addon.Event) {
	b.Publish(constant.Event, event)
}

func (b *Bus) PublishDeviceAdded(device *addon.Device) {
	b.Publish(constant.DeviceAdded, device)
}

func (b *Bus) PublishDeviceRemoved(deviceId string) {
	b.Publish(constant.DeviceRemoved, deviceId)
}

func (b *Bus) PublishThingConnected(thingId string, connected bool) {
	b.Publish(thingId+"."+constant.Connected, connected)
}

func (b *Bus) PublishThingRemoved(thingId string) {
	b.Publish(thingId + "." + constant.Removed)
}

func (b *Bus) AddDeviceRemovedSubscription(fn func(deviceId string)) func() {
	b.subscribe(constant.DeviceRemoved, fn)
	return func() {
		b.unsubscribe(constant.DeviceRemoved, fn)
	}
}

func (b *Bus) AddDeviceAddedSubscription(fn func(device *addon.Device)) func() {
	b.subscribe(constant.DeviceAdded, fn)
	return func() {
		b.unsubscribe(constant.DeviceAdded, fn)
	}
}

func (b *Bus) AddThingAddedSubscription(f func(thing *container.Thing)) func() {
	b.subscribe(constant.ThingAdded, f)
	return func() {
		b.unsubscribe(constant.ThingAdded, f)
	}
}

func (b *Bus) AddRemovedSubscription(thingId string, fn func()) func() {
	b.subscribe(thingId+"."+constant.Removed, fn)
	return func() {
		b.unsubscribe(thingId+"."+constant.Removed, fn)
	}
}

func (b *Bus) AddConnectedSubscription(thingId string, fn func(b bool)) func() {
	b.subscribe(thingId+"."+constant.Connected, fn)
	return func() {
		b.unsubscribe(thingId+"."+constant.Connected, fn)
	}
}

func (b *Bus) AddModifiedSubscription(thingId string, fn func()) func() {
	b.subscribe(thingId+"."+constant.Modified, fn)
	return func() {
		b.unsubscribe(thingId+"."+constant.Modified, fn)
	}
}

func (b *Bus) AddPropertyChangedSubscription(thingId string, fn func(p *addon.Property)) func() {
	b.subscribe(thingId+"."+constant.PropertyChanged, fn)
	return func() {
		b.unsubscribe(thingId+"."+constant.PropertyChanged, fn)
	}
}

func (b *Bus) AddActionStatusSubscription(f func(action *addon.Action)) func() {
	b.subscribe(constant.ActionStatus, f)
	return func() {
		b.unsubscribe(constant.ActionStatus, f)
	}
}

func (b *Bus) AddThingEventSubscription(f func(event *addon.Event)) func() {
	b.subscribe(constant.Event, f)
	return func() {
		b.unsubscribe(constant.Event, f)
	}
}

func (b *Bus) subscribe(topic string, fn interface{}) {
	b.logger.Debugf("subscribe topic:[%s]", topic)
	err := b.Bus.Subscribe(topic, fn)
	if err != nil {
		b.logger.Error("topic:%s subscribe err :%s", topic, err.Error())
	}
}

func (b *Bus) unsubscribe(topic string, fn interface{}) {
	b.logger.Debugf("unsubscribe topic:[%s]", topic)
	err := b.Bus.Unsubscribe(topic, fn)
	if err != nil {
		b.logger.Error("topic: %s unsubscribe err: %s", topic, err.Error())
	}
}

func (b *Bus) publish(topic string, args ...interface{}) {
	b.logger.Debugf("publish topic:[%s]", topic)
	if !b.Bus.HasCallback(topic) {
		return
	}
	b.Bus.Publish(topic, args...)
}

func (b *Bus) subscribeOnce(topic string, fn interface{}) {
	b.logger.Debugf("subscribeOnce topic:[%s]", topic)
	err := b.Bus.SubscribeOnce(topic, fn)
	if err != nil {
		b.logger.Error("topic: %s subscribe once err: %s", topic, err.Error())
	}
}

func (b *Bus) subscribeAsync(topic string, fn interface{}) {
	b.logger.Debugf("subscribeAsync topic:[%s]", topic)
	err := b.Bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		b.logger.Error("topic: %s subscribe async err: %s", topic, err.Error())
	}
}

func (b *Bus) waitAsync() {
	b.Bus.WaitAsync()
}
