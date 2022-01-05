package yeelight

import "github.com/galenliu/gateway/pkg/addon/properties"

type YeelightProperty struct {
	*properties.Property
}

func NewYeelightProperty(device properties.DeviceProxy, description properties.PropertyDescription) *YeelightProperty {
	return &YeelightProperty{properties.NewProperty(device, description)}
}
