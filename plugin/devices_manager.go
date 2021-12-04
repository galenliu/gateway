package plugin

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"time"
)

const SetPropertyTimeOut = 3000

func (m *Manager) SetPropertyValue(deviceId, propName string, newValue interface{}) (interface{}, error) {
	device := m.getDevice(deviceId)
	if device == nil {
		return nil, fmt.Errorf("device:%s not found", deviceId)
	}
	_, ok := device.GetProperty(propName)
	if !ok {
		return nil, fmt.Errorf("property:%s not found", propName)
	}
	device.setPropertyValue(propName, newValue)
	var valueChan = make(chan interface{})
	unsubscribeFunc := m.bus.Sub(topic.DevicePropertyChanged, func(deviceId string, p *addon.PropertyDescription) {
		if deviceId != device.GetId() {
			if p.Name == propName {
				valueChan <- p.Value
			}
		}
	})
	defer func() {
		unsubscribeFunc()
		close(valueChan)
	}()
	timeOut := time.After(SetPropertyTimeOut * time.Millisecond)
	for {
		select {
		case <-timeOut:
			return nil, fmt.Errorf("time out")
		case p := <-valueChan:
			return p, nil
		}
	}
}

func (m *Manager) GetPropertyValue(thingId, propName string) (interface{}, error) {
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

func (m *Manager) GetPropertiesValue(deviceId string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
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
	m.devices.Range(func(key, value interface{}) bool {
		device, ok := value.(*Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return
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
