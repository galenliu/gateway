package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

func (m *Manager) SetPropertyValue(ctx context.Context, deviceId, propName string, newValue any) (any, error) {
	device := m.getDevice(deviceId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", deviceId)
	}
	_, ok := device.GetProperty(propName)
	if !ok {
		return nil, fmt.Errorf("property:%s not found", propName)
	}
	return device.setPropertyValue(ctx, propName, newValue)
}

func (m *Manager) GetPropertyValue(thingId, propName string) (any, error) {
	device := m.getDevice(thingId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", thingId)
	}
	p, ok := device.GetProperty(propName)
	if !ok {
		return nil, fmt.Errorf("property:%s not found", propName)
	}
	return p.GetValue(), nil
}

func (m *Manager) GetPropertiesValue(deviceId string) (map[string]any, error) {
	data := make(map[string]any)
	device := m.getDevice(deviceId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", deviceId)
	}
	propMap := device.GetProperties()
	for n, p := range propMap {
		data[n] = p.GetValue()
	}
	return data, nil
}

func (m *Manager) GetMapOfDevices() (devices map[string]*addon.Device) {
	devs := m.GetDevices()
	var devicesMap = make(map[string]*addon.Device)
	if devs != nil {
		for _, dev := range devs {
			devicesMap[dev.GetId()] = dev.Device
		}
		return devicesMap
	}
	return
}

func (m *Manager) GetDevices() (devices []*Device) {
	m.devices.Range(func(key, value any) bool {
		device, ok := value.(*Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return
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
