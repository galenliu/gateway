package adapter

import "sync"

type Device interface {
	GetId() string
}

type AdapterProxy interface {
	GetId() string
	GetName() string
}

type Adapter struct {
	Id      string
	Name    string
	devices sync.Map
}

func (a *Adapter) AddDevice(dev Device) {
	a.devices.Store(dev.GetId(), dev)
}

func (a *Adapter) RemoveDevice(id string) {
	a.devices.Delete(id)
}

func (a *Adapter) GetDevice(id string) Device {
	i, ok := a.devices.Load(id)
	if ok {
		dev, ok := i.(Device)
		if ok {
			return dev
		}
	}
	return nil
}

func (a *Adapter) GetId() string {
	return a.Id
}

func (a *Adapter) GetName() string {
	return a.Name
}
