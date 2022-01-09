package yeelight

import (
	"github.com/akominch/yeelight"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

type YeelightDevice struct {
	*proxy.Device
	*yeelight.YeelightParams
	*yeelight.Bulb
}

func NewYeelightBulb(adapter addon.AdapterProxy, bulb *yeelight.Bulb, params *yeelight.YeelightParams) *YeelightDevice {
	yeeDevice := &YeelightDevice{
		Device:         proxy.NewDevice(adapter, []string{"Light", "OnOffSwitch"}, params.Name, "yeelight"+params.Name),
		YeelightParams: params,
		Bulb:           bulb,
	}
	for _, method := range params.Support {
		switch method {
		case "set_power":
			var atType = "OnOffProperty"
			prop := NewYeelightProperty(yeeDevice, addon.PropertyDescription{
				Name:        nil,
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
			yeeDevice.AddProperty("on", prop)
		case "set_bright":
			var min float64 = 0
			var max float64 = 100
			var atType = "LevelProperty"
			prop := NewYeelightProperty(yeeDevice, addon.PropertyDescription{
				Name:        nil,
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
			yeeDevice.AddProperty("bright", prop)
		case "set_rgb":
			var min float64 = 0
			var max float64 = 100
			var atType = "LevelProperty"
			prop := NewYeelightProperty(yeeDevice, addon.PropertyDescription{
				Name:        nil,
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
			yeeDevice.AddProperty("on", prop)
		default:
			continue
		}
	}
	return yeeDevice
}

func (d YeelightDevice) SetCredentials(username, password string) error {
	return nil
}
