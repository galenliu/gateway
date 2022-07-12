package addon

import (
	"context"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/errors"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

// SetPropertyValue 返回Fiber.NewError
func (m *Manager) SetPropertyValue(ctx context.Context, thingId, propName string, newValue any) (any, error) {
	device := m.getDevice(thingId)
	if device == nil {
		m.Publish(topic.DeviceConnected, topic.DeviceConnectedMessage{
			DeviceId:  thingId,
			Connected: false,
		})
		return nil, errors.NotFoundError("device %s not found", thingId)
	}
	p := device.GetProperty(propName)
	if p == nil {
		return nil, errors.NotFoundError("property %s not found", propName)
	}
	return device.setPropertyValue(ctx, propName, newValue)
}

func (m *Manager) GetPropertyValue(thingId, propName string) (any, error) {
	device := m.getDevice(thingId)
	if device == nil {
		return nil, errors.NotFoundError("device %s not found", thingId)
	}
	p := device.GetProperty(propName)
	if p == nil {
		return nil, errors.NotFoundError("property %s not found", propName)
	}
	return p.GetCachedValue(), nil
}

func (m *Manager) GetPropertiesValue(thingId string) (map[string]any, error) {
	data := make(map[string]any)
	device := m.getDevice(thingId)
	if device == nil {
		return nil, errors.NotFoundError("device:%s not found", thingId)
	}
	propMap := device.GetProperties()
	for n, p := range propMap {
		data[n] = p.GetCachedValue()
	}
	return data, nil
}

func (m *Manager) GetMapOfDevices() map[string]*devices.Device {
	devs := m.GetDevices()
	var devicesMap = make(map[string]*devices.Device, 0)
	if devs != nil {
		for _, dev := range devs {
			devicesMap[dev.GetId()] = dev.Device
		}
		return devicesMap
	}
	return devicesMap
}

func (m *Manager) GetDevices() (devices []*device) {
	devices = make([]*device, 0)
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
		return nil, errors.NotFoundError("device:%s not found", thingId)
	}
	return device.getAdapter().setDevicePin(ctx, thingId, pin)
}

func (m *Manager) SetCredentials(ctx context.Context, thingId, username, password string) (*messages.Device, error) {
	device := m.getDevice(thingId)
	if device == nil {
		return nil, errors.NotFoundError("device:%s not found", thingId)
	}
	return device.getAdapter().setDeviceCredentials(ctx, thingId, username, password)
}
