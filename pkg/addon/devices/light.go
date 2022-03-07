package devices

import (
	"fmt"
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
	OnOff            properties.BooleanEntity
	Brightness       properties.IntegerEntity
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
			light.OnOff = p.(properties.BooleanEntity)
		case properties.TypeBrightnessProperty:
			light.Brightness = p.(properties.IntegerEntity)
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

func (light *Light) TurnOn() error {
	return light.OnOff.TurnOn()
}

func (light *Light) TurnOff() error {
	return light.OnOff.TurnOff()
}

func (light *Light) Toggle() error {
	if light.OnOff.IsOn() {
		return light.TurnOn()
	} else {
		return light.TurnOff()
	}
}

func (light *Light) SetBrightness(brightness int) error {
	if light.Brightness != nil {
		return light.Brightness.SetValue(properties.Integer(brightness))
	}
	return fmt.Errorf("brightness not set for light")
}

func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}
