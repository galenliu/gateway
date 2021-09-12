package plugin

import (
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
)

type Eventbus struct {
	bus bus.Controller
}

func NewEventBus(bus bus.Controller) *Eventbus {
	b := &Eventbus{}
	b.bus = bus
	return b
}

func (e Eventbus) PublishConnected(deviceId string, connected bool) {
	e.bus.Publish(constant.CONNECTED, deviceId, connected)
}

func (e Eventbus) PublishActionStatus(action *rpc.DeviceActionStatusNotificationMessage_Data) {
	e.bus.Publish(constant.ActionStatus, action)
}

func (e Eventbus) SubscribeActionStatus(f func(action *rpc.ActionDescription)) {
	e.bus.Subscribe(constant.ActionStatus, f)
}

func (e Eventbus) UnsubscribeActionStatus(f func(action *rpc.ActionDescription)) {
	e.bus.Unsubscribe(constant.ActionStatus, f)
}

func (e Eventbus) PublishPropertyChanged(property *rpc.DevicePropertyChangedNotificationMessage_Data) {
	e.bus.Publish(constant.PropertyChanged, property)
}

func (e Eventbus) SubscribePropertyChanged(f func(property *rpc.DevicePropertyChangedNotificationMessage_Data)) {
	e.bus.Subscribe(constant.PropertyChanged, f)
}

func (e Eventbus) UnsubscribePropertyChanged(f func(property *rpc.DevicePropertyChangedNotificationMessage_Data)) {
	e.bus.Unsubscribe(constant.PropertyChanged, f)
}
