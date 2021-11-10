package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
)

func (m *Manager) SetPropertyValue(deviceId, propName string, newValue interface{}) (interface{}, error) {

	//addon, ok := m.devices[deviceId]
	//if !ok {
	//	return nil, fmt.Errorf("addon not found")
	//}
	//prop := addon.GetProperty(propName)
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
	//	Id := json.Get(data, "deviceId").ToString()
	//	name := json.Get(data, "name").ToString()
	//	value := json.Get(data, "value").GetInterface()
	//	if Id == deviceId && name == propName {
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
	//addon, ok := m.devices[deviceId]
	//if !ok {
	//	return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
	//}
	//prop := addon.GetProperty(propName)
	//if prop == nil {
	//	return nil, fmt.Errorf("propName(%s)invaild", propName)
	//}
	//return prop.GetValue(), nil
	return nil, nil
}

//func (m *Manager)GetPropertiesValue(deviceId string)(map[string]interface{},error){
//	addon, ok := m.devices[deviceId]
//	if !ok {
//		return nil, fmt.Errorf("deviceId (%s)invaild", deviceId)
//	}
//	addon.GetPropertyValue()
//
//}

func (m *Manager) GetDevice(deviceId string) *Device {
	return nil
}

func (m *Manager) GetDeviceMaps() (devices map[string]*addon.Device) {
	devs := m.getDevices()
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
	devs := m.getDevices()
	if devs != nil {
		for _, dev := range devs {
			devices = append(devices, dev)
		}
	}
	return devices
}

func (m *Manager) RemoveDevice(deviceId string) error {
	//
	//addon := m.getDevice(deviceId)
	//adapter := m.getAdapter(addon.AdapterId)
	//if adapter != nil {
	//	adapter.removeThing(addon)
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
		return fmt.Errorf("con not finid addon:" + thingId)
	}
	//err := addon.SetPin(pin)
	//if err != nil {
	//	return err
	//}
	return nil
}
