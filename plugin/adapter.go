package plugin

import (
	"addon"
	"fmt"
)

type Manager interface {
	handleDeviceRemoved(device *addon.Device)
	handleDeviceAdded(device *addon.Device)
}

type Adapter struct {
	ID          string `json:"adapterId"`
	Name        string `json:"name"`
	PackageName string `json:"packageName"`
	manager     Manager
	Devices     map[string]*addon.Device
	IsPairing   bool
}

func NewAdapter(adapterId, name, packageName string) *Adapter {
	adapter := &Adapter{}
	adapter.PackageName = packageName
	adapter.Name = name
	adapter.ID = adapterId
	adapter.Devices = make(map[string]*addon.Device, 10)
	adapter.IsPairing = false
	return adapter
}

func (adapter *Adapter) handleDeviceAdded(device *addon.Device) {
	if device == nil {
		return
	}
	device.AdapterId = adapter.ID
	adapter.Devices[device.ID] = device
	adapter.manager.handleDeviceAdded(device)
}

func (adapter *Adapter) handleDeviceRemoved(device *addon.Device) {
	delete(adapter.Devices, device.ID)
	adapter.manager.handleDeviceAdded(device)
}

func (adapter *Adapter) GetAdapterId() string {
	return adapter.ID
}

func (adapter *Adapter) GetPacketName() string {
	return adapter.PackageName
}

func (adapter *Adapter) GetManger() Manager {
	return adapter.manager
}

func (adapter *Adapter) FindDevice(deviceId string) (*addon.Device, error) {
	device, ok := adapter.Devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("devices id:(%s) invaild", deviceId)
	}
	return device, nil
}

func (adapter *Adapter) GetDevice(deviceId string) *addon.Device {
	device, ok := adapter.Devices[deviceId]
	if !ok {
		return nil
	}
	return device
}
