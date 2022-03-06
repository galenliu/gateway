package devices

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
	"strconv"
	"strings"
)

type LightHandler interface {
	TurnOn()
	TurnOff()
	SetBrightness(brightness int)
}

type Light struct {
	*Device
	OnOff            properties.Entity
	Brightness       properties.Entity
	ColorMode        properties.Entity
	Color            properties.Entity
	ColorTemperature properties.Entity
}

func NewLightBulb(id string, opts ...Option) *Light {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityLight, CapabilityOnOffSwitch},
	}, opts...)
	if device == nil {
		return nil
	}
	return &Light{Device: device}
}

func (light *Light) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeOnOffProperty:
			light.OnOff = p
		case properties.TypeBrightnessProperty:
			light.Brightness = p
		case properties.TypeColorModeProperty:
			light.ColorMode = p
		case properties.TypeColorProperty:
			light.Color = p
		case properties.TypeColorTemperatureProperty:
			light.ColorTemperature = p
		}
		light.Device.AddProperty(p)
	}
}

func (light *Light) TurnOn() {
	//light.On.SetValue(true)
}

func (light *Light) TurnOff() {
	//light.On.SetValue(false)
}

func (light *Light) Toggle() {
	//if light.On.Value == true {
	//	light.TurnOff()
	//} else {
	//	light.TurnOn()
	//}
}

func (light *Light) SetBrightness(brightness int) {
	//light.Device.get
}

func (light *Light) propertyValueUpdate(propName string, newValue any) {

}

func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}
