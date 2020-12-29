package event

import (
	"gateway/addons"
	"gateway/app/models"
	"sync"
)

type OnPropertyValueChangedFunc func(value interface{})
type OnDiscoverNewDeviceFunc func(addons.Device)
type OnActionStatusFunc func()
type RemoveFunc func()

var once sync.Once

var Bus *EventBus
var eid int

type EventBus struct {
	OnPropertyValueChangedListener map[string]map[string]map[int]OnPropertyValueChangedFunc
	OnDiscoverNewDeviceListener    map[int]OnDiscoverNewDeviceFunc
}

func InitEventBus() {
	once.Do(
		func() {
			eid = 0
			Bus = &EventBus{
				OnPropertyValueChangedListener: make(map[string]map[string]map[int]OnPropertyValueChangedFunc, 20),
				OnDiscoverNewDeviceListener:    make(map[int]OnDiscoverNewDeviceFunc, 20),
			}
		},
	)

}

func FirePropertyValueChanged(deviceId, PropName string, value interface{}) {
	for _, e := range Bus.OnPropertyValueChangedListener[deviceId][PropName] {
		e(value)
	}
}

func ListenPropertyValueChanged(deviceId, PropName string, f OnPropertyValueChangedFunc) func() {
	eid++
	var events = Bus.OnPropertyValueChangedListener[deviceId][PropName]
	events[eid] = f
	removeFunc := func() {
		delete(events, eid)
	}
	return removeFunc
}

func ListenDiscoverNewDevice(f OnDiscoverNewDeviceFunc) func() {
	eid++
	Bus.OnDiscoverNewDeviceListener[eid] = f
	removeFunc := func() {
		delete(Bus.OnDiscoverNewDeviceListener, eid)
	}
	return removeFunc
}

func FireDiscoverNewDevice(device *addons.Device) {
	for _, f := range Bus.OnDiscoverNewDeviceListener {
		f(*device)
	}
}

func ListenAction() {

}

func FireAction(action *models.Action) {

}
