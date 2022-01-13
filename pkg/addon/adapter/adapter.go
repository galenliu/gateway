package adapter

import "sync"

type Device interface {
	GetId() string
}

type Adapter struct {
	id          string
	packageName string
	devices     sync.Map
	name        string
}

func NewAdapter(id, packageName string) *Adapter {
	return &Adapter{
		id:          id,
		packageName: packageName,
		devices:     sync.Map{},
	}
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

func (a *Adapter) GetId() string {
	return a.id
}

func (a *Adapter) GetName() string {
	return a.name
}

func (a *Adapter) GetPackageName() string {
	return a.packageName
}
