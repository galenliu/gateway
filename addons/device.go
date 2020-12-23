package addons

import _ "github.com/asaskevich/govalidator"

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

	Context = "https://www.w3.org/2019/wot/td/v1"
)

type Device struct {
	AtContext    string   `json:"@context" valid:"url,required"`
	AtType       []string `json:"@type" valid:"url"`
	Title        string   `json:"title,required"`
	Titles       []string `json:"titles,omitempty"`
	Description  string   `json:"description,omitempty"`
	ID           string   `json:"id"`
	Descriptions []string `json:"descriptions,omitempty"`
	Version      string   `json:"version,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
	Support      string   `json:"support,omitempty"`
	Properties   map[string]*Property
}

func (proxy *Device) FindProperty(propName string) *Property {
	prop := proxy.Properties[propName]
	return prop
}

func (proxy *Device) GetProperty(propName string) interface{} {
	prop := proxy.FindProperty(propName)
	return prop.Value
}
