package plugin

import (
	"fmt"
	addon "github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"time"
)

func (m *Manager) SetPropertyValue(deviceId, propName string, newValue interface{}) (interface{}, error) {

	go func() {
		err := m.handleSetProperty(deviceId, propName, newValue)
		if err != nil {
			m.logger.Error(err.Error())
		}
	}()
	closeChan := make(chan struct{})
	propChan := make(chan interface{})
	time.AfterFunc(3*time.Second, func() {
		closeChan <- struct{}{}
	})
	changed := func(data []byte) {
		id := json.Get(data, "deviceId").ToString()
		name := json.Get(data, "name").ToString()
		value := json.Get(data, "value").GetInterface()
		if id == deviceId && name == propName {
			propChan <- value
		}
	}
	go m.bus.Subscribe(util.PropertyChanged, changed)
	defer m.bus.Unsubscribe(util.PropertyChanged, changed)
	for {
		select {
		case v := <-propChan:
			return v, nil
		case <-closeChan:
			m.logger.Error("set property(name: %s value: %s) timeout", propName, newValue)
			return nil, fmt.Errorf("timeout")
		}
	}
}

func (m *Manager) GetDevice(deviceId string) addon.IDevice {
	device, ok := m.devices[deviceId]
	if !ok {
		return nil
	}
	return device
}

func (m *Manager) GetDevices() (device []addon.IDevice) {
	for _, dev := range m.devices {
		device = append(device, dev)
	}
	return
}

func (m *Manager) RemoveDevice(deviceId string) error {

	device := m.getDevice(deviceId)
	adapter := m.getAdapter(device.GetAdapterId())
	if adapter != nil {
		adapter.removeThing(device)
		return nil
	}
	return fmt.Errorf("can not find thing")
}

func (m *Manager) GetPropertyValue(deviceId, propName string) (interface{}, error) {
	device, ok := m.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
	}
	prop := device.GetProperty(propName)
	if prop == nil {
		return nil, fmt.Errorf("propName(%s)invaild", propName)
	}
	return prop.GetValue(), nil
}