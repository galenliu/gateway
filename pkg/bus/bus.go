package bus

import (
	"fmt"
	"gateway/log"
	"github.com/asaskevich/EventBus"
)

var bus EventBus.Bus

func initBus() {
	bus = EventBus.New()
}

func Subscribe(topic string, fn interface{}) error {

	if bus == nil {
		initBus()
	}
	return bus.Subscribe(topic, fn)
}

func Unsubscribe(topic string, fn interface{}) error {
	if bus == nil {
		initBus()
	}
	return bus.Unsubscribe(topic, fn)
}

func Publish(topic string, args ...interface{}) {
	log.Info("publish topic: " + topic)
	if bus == nil {
		initBus()
	}
	log.Info(fmt.Sprintf(topic+" has callback %v", bus.HasCallback(topic)))
	if !bus.HasCallback(topic) {
		return
	}
	bus.Publish(topic, args...)
}

func HasCallBack(topic string) bool {
	return bus.HasCallback(topic)
}

func SubscribeOnce(topic string, fn interface{}) {
	if bus == nil {
		initBus()
	}
	_ = bus.SubscribeOnce(topic, fn)
}

func SubscribeAsync(topic string, fn interface{}) error {
	if bus == nil {
		initBus()
	}
	return bus.SubscribeAsync(topic, fn, false)
}

func WaitAsync() {
	if bus == nil {
		initBus()
	}
	bus.WaitAsync()
}
