package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

func (m *Manager) SetPropertyValue(ctx context.Context, deviceId, propName string, newValue any) (any, error) {
	device := m.getDevice(deviceId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", deviceId)
	}
	p := device.GetPropertyEntity(propName)
	if p == nil {
		return nil, fmt.Errorf("property:%s not found", propName)
	}
	return device.setPropertyValue(ctx, propName, newValue)
}

func (m *Manager) GetPropertyValue(thingId, propName string) (any, error) {
	device := m.getDevice(thingId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", thingId)
	}
	p := device.GetPropertyEntity(propName)
	if p == nil {
		return nil, fmt.Errorf("property:%s not found", propName)
	}
	return p.GetCachedValue(), nil
}

func (m *Manager) GetPropertiesValue(deviceId string) (map[string]any, error) {
	data := make(map[string]any)
	device := m.getDevice(deviceId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", deviceId)
	}
	propMap := device.GetProperties()
	for n, p := range propMap {
		data[n] = p.GetCachedValue()
	}
	return data, nil
}

func (m *Manager) GetMapOfDevices() map[string]*devices.Device {
	devs := m.GetDevices()
	var devicesMap = make(map[string]*devices.Device)
	if devs != nil {
		for _, dev := range devs {
			devicesMap[dev.GetId()] = dev.Device
		}
		return devicesMap
	}
	return nil
}

func (m *Manager) GetDevices() (devices []*device) {
	devices = make([]*device, 1)
	for _, s := range m.Manager.GetDevices() {
		device, ok := s.(*device)
		if ok {
			devices = append(devices, device)
		}
	}
	return devices
}

func (m *Manager) SetPIN(ctx context.Context, thingId string, pin string) (*messages.Device, error) {
	device := m.getDevice(thingId)
	if device == nil {
		return nil, fmt.Errorf("con not finid addon:" + thingId)
	}
	return device.getAdapter().setDevicePin(ctx, thingId, pin)
}

func (m *Manager) SetCredentials(ctx context.Context, thingId, username, password string) (*messages.Device, error) {
	device := m.getDevice(thingId)
	if device == nil {
		return nil, fmt.Errorf("con not finid addon:" + thingId)
	}
	return device.getAdapter().setDeviceCredentials(ctx, thingId, username, password)
}
