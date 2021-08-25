package bus

import (
	bus "github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/logging"
)

type Bus struct {
	bus.Bus
	logger logging.Logger
}

func NewBus(log logging.Logger) (*Bus, error) {
	b := &Bus{}
	b.logger = log
	b.Bus = bus.New()
	return b, nil
}

func (b *Bus) Subscribe(topic string, fn interface{}) {
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
	if !b.Bus.HasCallback(topic) {
		b.logger.Infof("topic[%s] has not callback", topic)
		return
	}
	b.logger.Infof("bus publish topic:[%s]", topic)
	b.Bus.Publish(topic, args...)
}

func (b *Bus) SubscribeOnce(topic string, fn interface{}) {
	err := b.Bus.SubscribeOnce(topic, fn)
	if err != nil {
		b.logger.Error("topic: %s subscribe once err: %s", topic, err.Error())
	}
}

func (b *Bus) SubscribeAsync(topic string, fn interface{}) {
	err := b.Bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		b.logger.Error("topic: %s subscribe async err: %s", topic, err.Error())
	}
}

func (b *Bus) WaitAsync() {
	b.Bus.WaitAsync()
}
