package manager

import "sync"

type Device interface {
	GetId() string
}

type Adapter interface {
	GetId() string
}

type Integration interface {
	GetId() string
}

func NewManager() *Manager {
	return &Manager{}
}

type Manager struct {
	devices      sync.Map
	adapters     sync.Map
	integrations sync.Map
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

func (m *Manager) AddIntegration(ig Integration) {
	m.integrations.Store(ig.GetId(), ig)
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

func (m *Manager) GetIntegration(id string) Integration {
	v, ok := m.integrations.Load(id)
	if ok {
		v, ok := v.(Integration)
		if ok {
			return v
		}
	}
	return nil
}

func (m *Manager) GetIntegrations() []Integration {
	integrations := make([]Integration, 1)
	m.integrations.Range(func(key, value any) bool {
		com, ok := value.(Integration)
		if ok {
			integrations = append(integrations, com)
		}
		return true
	})
	return integrations
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
