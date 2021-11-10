package container

import (
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
)

type EventBus struct {
	bus     bus.Controller
	thingId string
}

func NewEventBus(bus bus.Controller) *EventBus {
	return &EventBus{bus: bus}
}

func (bus EventBus) PublishConnected(connected bool) {
	topic := bus.thingId + "." + constant.CONNECTED
	bus.bus.Publish(topic, connected)
}

func (bus EventBus) PublishRemoved() {
	topic := bus.thingId + "." + constant.ThingRemoved
	bus.bus.Publish(topic)
}
