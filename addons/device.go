package addons

import (
	"fmt"
	_ "github.com/asaskevich/govalidator"
)

const (
	Alarm                    = "Alarm"
	AirQualitySensor         = "AirQualitySensor"
	BarometricPressureSensor = "BarometricPressureSensor"
	BinarySensor             = "BinarySensor"
	Camera                   = "Camera"
	ColorControl             = "ColorControl"
	ColorSensor              = "ColorSensor"
	DoorSensor               = "DoorSensor"
	EnergyMonitor            = "EnergyMonitor"
	HumiditySensor           = "HumiditySensor"
	LeakSensor               = "LeakSensor"
	Light                    = "Light"
	Lock                     = "Lock"
	MotionSensor             = "MotionSensor"
	MultiLevelSensor         = "MultiLevelSensor"
	MultiLevelSwitch         = "MultiLevelSwitch"
	OnOffSwitch              = "OnOffSwitch"
	SmartPlug                = "SmartPlug"
	SmokeSensor              = "SmokeSensor"
	TemperatureSensor        = "TemperatureSensor"
	Thermostat               = "Thermostat"
	VideoCamera              = "VideoCamera"

	Context = "https://webthings.io/schemas"
)

type Device struct {
	adapter Adapter

	AtContext   string                    `json:"@context" valid:"url,required"`
	Title       string                    `json:"title,required"`
	AtType      []string                  `json:"@type"`
	Description string                    `json:"description,omitempty"`
	ID          string                    `json:"id"`
	Properties  map[string]*Property `json:"properties"`
}

func NewDevice(id, title string) *Device {
	device := &Device{}
	device.ID = id
	device.Title = title
	device.Properties = make(map[string]*Property)
	return device
}

func (device *Device) AddProperties(props ...*Property) {
	for _, p := range props {
		device.Properties[p.Name] = p
	}
}


func (device *Device) AddTypes(types ...string) {
	for _, t := range types {
		device.AtType = append(device.AtType, t)
	}
}


func (device *Device) SetProperty(propName string, value interface{}) (interface{}, error) {
	prop, ok := device.Properties[propName]
	if !ok {
		return nil, fmt.Errorf("can find prop(%s) form device(%s)", device.ID, propName)
	}
	prop.setValue(value)
	return prop.Value, nil
}


