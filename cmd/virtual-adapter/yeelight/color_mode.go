package yeelight

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type YeeColorModeProperty struct {
	device *YeelightDevice
	*properties.ColorModeProperty
}

func NewColorMode(device *YeelightDevice, value string) *YeeColorModeProperty {
	return &YeeColorModeProperty{
		device,
		properties.NewColorModeProperty(colorMode(value), properties.WithReadOnly()),
	}
}

func colorMode(in string) properties.ColorModePropertyEnum {
	if in == "2" {
		return properties.ColorModePropertyEnumColor
	}
	return properties.ColorModePropertyEnumTemperature
}
