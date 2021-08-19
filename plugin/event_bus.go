package plugin

import "github.com/galenliu/gateway/pkg/constant"

type bus interface {
	Publish(string, ...interface{})
	Subscribe(string, interface{})
	Unsubscribe(string, interface{})
}

type Eventbus struct {
	bus bus
}

func NewEventBus(bus bus) *Eventbus {
	b := &Eventbus{}
	return b
}

func (e Eventbus) PublishPropertyChanged(data []byte) {
	e.bus.Publish(constant.PropertyChanged, data)
}

func (e Eventbus) SubscribePropertyChanged(f func(deviceId, propName string, data []byte)) {
	e.bus.Subscribe(constant.PropertyChanged, f)
}
