package bus

import (
	"github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/wot"
	"github.com/galenliu/gateway/wot/models"
)

type EventBusController interface {
	wot.ThingsEventBus
}

type bus struct {
	logger logging.Logger
	EventBus.Bus
}

func (b *bus) ListenCreateThing(f func(data []byte) error) {
	panic("implement me")
}

func (b *bus) ListenRemoveThing(f func(id string)) {
	panic("implement me")
}

func (b *bus) FireThingAdded(thing *models.Thing) {
	panic("implement me")
}

func (b *bus) FireThingRemoved(id string) {
	panic("implement me")
}

func NewEventBus(log logging.Logger) (EventBusController, error) {
	bus := &bus{}
	bus.logger = log
	bus.Bus = EventBus.New()
	return bus, nil
}

func (bus *bus) subscribe(topic string, fn interface{}) {
	err := bus.Subscribe(topic, fn)
	if err != nil {
		bus.logger.Error("topic:%s subscribe err :%s", topic, err.Error())
	}
}

func (bus *bus) unsubscribe(topic string, fn interface{}) {
	err := bus.Unsubscribe(topic, fn)
	if err != nil {
		bus.logger.Error("topic: %s unsubscribe err: %s", topic, err.Error())
	}
}

func (bus *bus) publish(topic string, args ...interface{}) {
	if !bus.HasCallback(topic) {
		bus.logger.Info("topic has not callback")
		return
	}
	bus.Publish(topic, args...)
}

func (bus *bus) subscribeOnce(topic string, fn interface{}) {
	err := bus.SubscribeOnce(topic, fn)
	if err != nil {
		bus.logger.Error("topic: %s subscribe once err: %s", topic, err.Error())
	}
}

func (bus *bus) subscribeAsync(topic string, fn interface{}) {
	err := bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		bus.logger.Error("topic: %s subscribe async err: %s", topic, err.Error())
	}
}

func (bus *bus) waitAsync() {
	bus.WaitAsync()
}

//type OnPropertyChangeFunc = func(property *addon.Property)
//
//func SubscribePropertyChanged(thingId, propName string, fn OnPropertyChangeFunc) error {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, thingId, propName)
//	if instance == nil {
//		initBus()
//	}
//	return instance.Subscribe(topic, fn)
//}
//
//func SubscribeOncePropertyChanged(thingId, propName string, fn OnPropertyChangeFunc) error {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, thingId, propName)
//
//	if instance == nil {
//		initBus()
//	}
//	return instance.SubscribeOnce(topic, fn)
//}
//
//func UnsubscribePropertyChanged(thingId, propName string, fn OnPropertyChangeFunc) error {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, thingId, propName)
//	if instance == nil {
//		initBus()
//	}
//	return instance.Unsubscribe(topic, fn)
//}
//
//func PublishPropertyChanged(prop *addon.Property) {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, prop.DeviceId, prop.name)
//	if instance == nil {
//		initBus()
//	}
//	log.Info(fmt.Sprintf(topic+" has callback %v", instance.HasCallback(topic)))
//	if !instance.HasCallback(topic) {
//		return
//	}
//	instance.Publish(topic, prop)
//}
