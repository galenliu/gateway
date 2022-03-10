package manager

import "sync"

type Device interface {
	GetId() string
}

type Adapter interface {
	GetId() string
}

func NewManager() *Manager {
	return &Manager{}
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

func (m *Manager) AddAdapter(a Adapter) {
	m.adapters.Store(a.GetId(), a)
}

func (m *Manager) RemoveAdapter(id string) {
	m.adapters.Delete(id)
}

func (m *Manager) GetAdapter(id string) Adapter {
	v, ok := m.adapters.Load(id)
	if ok {
		v, ok := v.(Adapter)
		if ok {
			return v
		}
	}
	return nil
}

func (m *Manager) GetAdapters() []Adapter {
	adapters := make([]Adapter, 1)
	m.adapters.Range(func(key, value any) bool {
		adp, ok := value.(Adapter)
		if ok {
			adapters = append(adapters, adp)
		}
		return true
	})
	return adapters
}

func (m *Manager) GetDevice(id string) Device {
	v, ok := m.devices.Load(id)
	if ok {
		v, ok := v.(Device)
		if ok {
			return v
		}
	}
	return nil
}

func (m *Manager) GetDevices() []Device {
	devices := make([]Device, 1)
	m.devices.Range(func(key, value any) bool {
		device, ok := value.(Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return nil
}
