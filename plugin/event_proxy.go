package plugin

import "github.com/galenliu/gateway-addon/events"

type Event struct {
	device *Device
	*events.Event
}

func NewEvent(device *Device, event *events.Event) *Event {
	return &Event{
		device: device,
		Event:  event,
	}
}
