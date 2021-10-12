package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway-addon/devices"
)

func (m *Manager) SetPropertyValue(deviceId, propName string, newValue interface{}) (interface{}, error) {

	//device, ok := m.devices[deviceId]
	//if !ok {
	//	return nil, fmt.Errorf("device not found")
	//}
	//prop := device.GetProperty(propName)
	//if prop == nil {
	//	return nil, fmt.Errorf("property not found")
	//}
	//
	//return prop.SetValue(newValue)

	//go func() {
	//	err := m.handleSetProperty(deviceId, propName, newValue)
	//	if err != nil {
	//		m.logger.Error(err.Error())
	//	}
	//}()
	//closeChan := make(chan struct{})
	//propChan := make(chan interface{})
	//time.AfterFunc(3*time.Second, func() {
	//	closeChan <- struct{}{}
	//})
	//changed := func(data []byte) {
	//	ID := json.Get(data, "deviceId").ToString()
	//	name := json.Get(data, "name").ToString()
	//	value := json.Get(data, "value").GetInterface()
	//	if ID == deviceId && name == propName {
	//		propChan <- value
	//	}
	//}
	//go m.bus.Subscribe(constant.PropertyChanged, changed)
	//defer m.bus.Unsubscribe(constant.PropertyChanged, changed)
	//for {
	//	select {
	//	case v := <-propChan:
	//		return v, nil
	//	case <-closeChan:
	//		m.logger.Error("set property(name: %s value: %s) timeout", propName, newValue)
	//		return nil, fmt.Errorf("timeout")
	//	}
	//}
	return nil, nil
}

func (m *Manager) GetPropertyValue(deviceId, propName string) (interface{}, error) {
	//device, ok := m.devices[deviceId]
	//if !ok {
	//	return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
	//}
	//prop := device.GetProperty(propName)
	//if prop == nil {
	//	return nil, fmt.Errorf("propName(%s)invaild", propName)
	//}
	//return prop.GetValue(), nil
	return nil, nil
}

//func (m *Manager)GetPropertiesValue(deviceId string)(map[string]interface{},error){
//	device, ok := m.devices[deviceId]
//	if !ok {
//		return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
//	}
//	device.GetPropertyValue()
//
//}

func (m *Manager) GetDevice(deviceId string) *Device {
	return nil
}

func (m *Manager) GetDevicesBytes() (devices map[string][]byte) {
	devs := m.getDevices()
	var devicesMap = make(map[string][]byte)
	if devs != nil {
		for _, dev := range devs {
			bt, err := json.Marshal(dev)
			if err != nil {
				devicesMap[dev.ID] = bt
			}
		}
		return devicesMap
	}
	return
}

func (m *Manager) GetDevices() (devices []*devices.Device) {
	devs := m.getDevices()
	if devs != nil {
		for _, dev := range devs {
			devices = append(devices, dev.Device)
		}
	}
	return devices
}

func (m *Manager) GetDevicesMaps() (dms map[string]*devices.Device) {
	devs := m.getDevices()
	if devs != nil {
		dms = make(map[string]*devices.Device)
		for _, dev := range devs {
			dms[dev.ID] = dev.Device
		}
	}
	return
}

func (m *Manager) RemoveDevice(deviceId string) error {
	//
	//device := m.getDevice(deviceId)
	//adapter := m.getAdapter(device.AdapterId)
	//if adapter != nil {
	//	adapter.removeThing(device)
	//	return nil
	//}
	return fmt.Errorf("can not find thing")
}

func (m *Manager) CancelRemoveThing(deviceId string) {
	device := m.getDevice(deviceId)
	if device == nil {
		return
	}

	if device.adapter != nil {
		device.adapter.cancelRemoveThing(deviceId)
	}
}

func (m *Manager) SetPIN(thingId string, pin interface{}) error {
	device := m.getDevice(thingId)
	if device == nil {
		return fmt.Errorf("con not finid device:" + thingId)
	}
	//err := device.SetPin(pin)
	//if err != nil {
	//	return err
	//}
	return nil
}
