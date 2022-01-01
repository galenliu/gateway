package addon

import (
	"github.com/galenliu/gateway/pkg/addon/gateway-addon/devices"
	"github.com/galenliu/gateway/pkg/addon/gateway-addon/properties"
)

type DeviceProxy struct {
	*devices.Device
}

func (d DeviceProxy) AddProperty(property *properties.Property) {
	d.Properties[property.Name] = property
}
