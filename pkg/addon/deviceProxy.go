package addon

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type DeviceProxy struct {
	*devices.Device
}

func (d DeviceProxy) AddProperty(property properties.Property) {
	d.Properties[property.GetName()] = property
}
