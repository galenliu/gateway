package yeelight

import (
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

type YeelightDevice struct {
	*proxy.Device
	*yeelight.Yeelight
}

func NewYeelightBulb(bulb *yeelight.Yeelight) *YeelightDevice {
	yeeDevice := &YeelightDevice{
		Device:   proxy.NewDevice([]string{"Light", "OnOffSwitch"}, bulb.GetAddr(), "yeelight"+bulb.GetAddr()),
		Yeelight: bulb,
	}
	for _, method := range bulb.GetSupports() {
		switch method {
		case "set_power":
			var atType = addon.OnOffProperty
			prop := NewYeelightProperty(bulb, properties.PropertyDescription{
				Name:        &on,
				AtType:      &atType,
				Title:       nil,
				Type:        addon.TypeBoolean,
				Unit:        nil,
				Description: nil,
				Minimum:     nil,
				Maximum:     nil,
				Enum:        nil,
				ReadOnly:    nil,
				MultipleOf:  nil,
				Links:       nil,
				Value:       nil,
			})
			yeeDevice.AddProperty(prop)
		case "set_bright":
			var min float64 = 0
			var max float64 = 100
			var atType = addon.LevelProperty
			prop := NewYeelightProperty(bulb, properties.PropertyDescription{
				Name:        &level,
				AtType:      &atType,
				Title:       nil,
				Type:        addon.TypeInteger,
				Unit:        nil,
				Description: nil,
				Minimum:     &min,
				Maximum:     &max,
				Enum:        nil,
				ReadOnly:    nil,
				MultipleOf:  nil,
				Links:       nil,
				Value:       nil,
			})
			yeeDevice.AddProperty(prop)
		case "set_rgb":
			var atType = addon.ColorProperty
			prop := NewYeelightProperty(bulb, properties.PropertyDescription{
				Name:        &color,
				AtType:      &atType,
				Title:       nil,
				Type:        addon.TypeString,
				Unit:        nil,
				Description: nil,
				Enum:        nil,
				ReadOnly:    nil,
				MultipleOf:  nil,
				Links:       nil,
				Value:       nil,
			})
			yeeDevice.AddProperty(prop)
		default:
			continue
		}
	}
	return yeeDevice
}

func (d YeelightDevice) SetCredentials(username, password string) error {
	return nil
}
