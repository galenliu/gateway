package adapter

import (
	"sync"
)

type Device interface {
	GetId() string
}

type Thing interface {
	GetId() string
}

type Entity interface {
	GetId() string
	GetDevice(id string) Device
	GetDevices() []Device
	RemoveDevice(deviceId string)
	GetAdapter() *Adapter
}

type Adapter struct {
	id      string
	devices sync.Map
}

func NewAdapter(id string) *Adapter {
	return &Adapter{
		id:      id,
		devices: sync.Map{},
	}
}

func (a *Adapter) GetAdapter() *Adapter {
	return a
}

func (a *Adapter) AddDevice(dev Device) {
	a.devices.Store(dev.GetId(), dev)
}

func (a *Adapter) RemoveDevice(id string) {
	a.devices.Delete(id)
}

func (a *Adapter) GetDevice(id string) Device {
	v, ok := a.devices.Load(id)
	if ok {
		v, ok := v.(Device)
		if ok {
			return v
		}
	}
	return nil
}

func (a *Adapter) GetDevices() []Device {
	devices := make([]Device, 0)
	a.devices.Range(func(key, value any) bool {
		device, ok := value.(Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return devices
}

func (a *Adapter) GetId() string {
	return a.id
}
