package yeelight

import (
	_type "github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type YeelightProperty struct {
	*properties.Property
}

func NewYeelightProperty(device _type.DeviceProxy, description _type.PropertyDescription) *YeelightProperty {
	return &YeelightProperty{properties.NewProperty(device, description)}
}
