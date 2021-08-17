package bus

import (
	bus "github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/logging"
)

type eventBus struct {
	bus    bus.Bus
	logger logging.Logger
}

type BaseBus interface {
	Subscribe(topic string, fn interface{})
	Unsubscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
	SubscribeOnce(topic string, fn interface{})
	SubscribeAsync(topic string, fn interface{})
	WaitAsync()
}

func NewEventBus(log logging.Logger) (eventBus, error) {
	b := eventBus{}
	b.logger = log
	b.bus = bus.New()
	return b, nil
}

func (eventBus eventBus) Subscribe(topic string, fn interface{}) {
	err := eventBus.bus.Subscribe(topic, fn)
	if err != nil {
		eventBus.logger.Error("topic:%s subscribe err :%s", topic, err.Error())
	}
}

func (eventBus eventBus) Unsubscribe(topic string, fn interface{}) {
	err := eventBus.bus.Unsubscribe(topic, fn)
	if err != nil {
		eventBus.logger.Error("topic: %s unsubscribe err: %s", topic, err.Error())
	}
}

func (eventBus eventBus) Publish(topic string, args ...interface{}) {
	if !eventBus.bus.HasCallback(topic) {
		eventBus.logger.Info("topic has not callback")
		return
	}
	eventBus.bus.Publish(topic, args...)
}

func (eventBus eventBus) SubscribeOnce(topic string, fn interface{}) {
	err := eventBus.bus.SubscribeOnce(topic, fn)
	if err != nil {
		eventBus.logger.Error("topic: %s subscribe once err: %s", topic, err.Error())
	}
}

func (eventBus eventBus) SubscribeAsync(topic string, fn interface{}) {
	err := eventBus.bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		eventBus.logger.Error("topic: %s subscribe async err: %s", topic, err.Error())
	}
}

func (eventBus eventBus) WaitAsync() {
	eventBus.bus.WaitAsync()
}
