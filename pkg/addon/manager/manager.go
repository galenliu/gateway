package manager

import "sync"

type Device interface {
	GetId() string
}

type Adapter interface {
	GetId() string
}

type Manager struct {
	devices  sync.Map
	adapters sync.Map
}

func (m *Manager) AddDevice(d Device) {
	m.devices.Store(d.GetId(), d)
}

func (m *Manager) RemoveDevice(id string) {
	m.devices.Delete(id)
}

func (m *Manager) GetDevice(id string) Device {
	i, ok := m.devices.Load(id)
	if ok {
		dev, ok := i.(Device)
		if ok {
			return dev
		}
	}
	return nil
}

func (m *Manager) AddAdapter(a Adapter) {
	m.devices.Store(a.GetId(), a)
}

func (m *Manager) RemoveAdapter(id string) {
	m.adapters.Delete(id)
}

func (m *Manager) GetAdapter(id string) Adapter {
	i, ok := m.adapters.Load(id)
	if ok {
		a, ok := i.(Adapter)
		if ok {
			return a
		}
	}
	return nil
}

func (m *Manager) GetAdapters() []Adapter {
	adapters := make([]Adapter, 1)
	m.adapters.Range(func(id any, v any) bool {
		a, ok := v.(Adapter)
		if ok {
			adapters = append(adapters, a)
		}
		return true
	})
	return adapters
}
