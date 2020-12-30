package event

import (
	addon "gitee.com/liu_guilin/gateway-addon-golang"
	"sync"
)

type AdapterProxy interface {
}

type PropertyProxy interface {
}
type OnPropertyChangedFunc func(value interface{})
type OnDiscoverNewDeviceFunc func(proxy addon.Device)
type OnActionStatusFunc func()
type RemoveFunc func()

var once sync.Once

var Bus *EventBus
var eid int

type EventBus struct {
	OnPropertyChangedListener   map[string]map[string]map[int]OnPropertyChangedFunc
	OnDiscoverNewDeviceListener map[int]OnDiscoverNewDeviceFunc
}

func InitEventBus() {
	once.Do(
		func() {
			eid = 0
			Bus = &EventBus{
				OnPropertyChangedListener:   make(map[string]map[string]map[int]OnPropertyChangedFunc, 20),
				OnDiscoverNewDeviceListener: make(map[int]OnDiscoverNewDeviceFunc, 20),
			}
		},
	)

}

func FireAdapterAdded(adapter interface{}) {

}

func FirePropertyChanged(property interface{}) {

}

func ListenPropertyValueChanged(deviceId, PropName string, f OnPropertyChangedFunc) func() {
	eid++
	var events = Bus.OnPropertyChangedListener[deviceId][PropName]
	events[eid] = f
	removeFunc := func() {
		delete(events, eid)
	}
	return removeFunc
}

func ListenDiscoverNewDevice(f OnDiscoverNewDeviceFunc) func() {
	//eid++
	//Bus.OnDiscoverNewDeviceListener[eid] = f
	removeFunc := func() {
		//delete(Bus.OnDiscoverNewDeviceListener, eid)
	}
	return removeFunc
}

func FireDiscoverNewDevice(device *addon.Device) {
	for _, f := range Bus.OnDiscoverNewDeviceListener {
		f(*device)
	}
}

func ListenAction() {

}

func FireAction() {

}
