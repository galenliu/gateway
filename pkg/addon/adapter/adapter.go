package adapter

import "sync"

type Device interface {
	GetId() string
}

type Thing interface {
	GetId() string
}

type Entity interface {
	GetId() string
	GetDeviceById(id string) Device
	GetDevices() []Device
}

type Adapter struct {
	id      string
	devices sync.Map
	things  sync.Map
}

func NewAdapter(id string) *Adapter {
	return &Adapter{
		id:      id,
		devices: sync.Map{},
	}
}

func (a *Adapter) AddDevice(dev Device) {
	a.devices.Store(dev.GetId(), dev)
}

func (a *Adapter) RemoveDevice(id string) {
	a.devices.Delete(id)
}

func (a *Adapter) GetDeviceById(id string) Device {
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
	devices := make([]Device, 1)
	a.devices.Range(func(key, value any) bool {
		device, ok := value.(Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return devices
}

func (a *Adapter) GetThings() []Thing {
	things := make([]Thing, 1)
	a.things.Range(func(key, value any) bool {
		device, ok := value.(Device)
		if ok {
			things = append(things, device)
		}
		return true
	})
	return things
}

func (a *Adapter) GetThingById(id string) Thing {
	v, ok := a.things.Load(id)
	if ok {
		v, ok := v.(Thing)
		if ok {
			return v
		}
	}
	return nil
}

func (a *Adapter) GetId() string {
	return a.id
}
