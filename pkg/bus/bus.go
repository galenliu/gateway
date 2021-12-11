package bus

import (
	bus "github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/bus/topic"
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

func (b *Bus) Sub(topic topic.Topic, fn interface{}) func() {
	go b.subscribe(topic.ToString(), fn)
	return func() {
		b.unsubscribe(topic.ToString(), fn)
	}
}

func (b *Bus) Pub(topic topic.Topic, args ...interface{}) {
	b.logger.Debugf("publish topic:[%s] Args:%#v", topic, args)
	go b.Bus.Publish(topic.ToString(), args...)
}

func (b *Bus) waitAsync() {
	b.Bus.WaitAsync()
}
