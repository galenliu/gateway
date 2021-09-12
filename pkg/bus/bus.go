package bus

import (
	bus "github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/logging"
)

type Controller interface {
	Subscribe(topic string, fn interface{})
	Unsubscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
	SubscribeOnce(topic string, fn interface{})
	SubscribeAsync(topic string, fn interface{})
}

type Bus struct {
	bus.Bus
	logger logging.Logger
}

func NewController(log logging.Logger) (Controller, error) {
	b := &Bus{}
	b.logger = log
	b.Bus = bus.New()
	return b, nil
}

func (b *Bus) Subscribe(topic string, fn interface{}) {
	b.logger.Debugf("subscribe topic:[%s]", topic)
	err := b.Bus.Subscribe(topic, fn)
	if err != nil {
		b.logger.Error("topic:%s subscribe err :%s", topic, err.Error())
	}
}

func (b *Bus) Unsubscribe(topic string, fn interface{}) {
	err := b.Bus.Unsubscribe(topic, fn)
	if err != nil {
		b.logger.Error("topic: %s unsubscribe err: %s", topic, err.Error())
	}
}

func (b *Bus) Publish(topic string, args ...interface{}) {
	b.logger.Debugf("publish topic:[%s]", topic)
	if !b.Bus.HasCallback(topic) {
		return
	}
	b.Bus.Publish(topic, args...)
}

func (b *Bus) SubscribeOnce(topic string, fn interface{}) {
	b.logger.Debugf("subscribeOnce topic:[%s]", topic)
	err := b.Bus.SubscribeOnce(topic, fn)
	if err != nil {
		b.logger.Error("topic: %s subscribe once err: %s", topic, err.Error())
	}
}

func (b *Bus) SubscribeAsync(topic string, fn interface{}) {
	b.logger.Debugf("SubscribeAsync topic:[%s]", topic)
	err := b.Bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		b.logger.Error("topic: %s subscribe async err: %s", topic, err.Error())
	}
}

func (b *Bus) WaitAsync() {
	b.Bus.WaitAsync()
}
