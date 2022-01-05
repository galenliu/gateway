package pkg

import (
	"github.com/galenliu/gateway/cmd/virtual-adapter/pkg/yeelight"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type VirtualProperty struct {
	*properties.Property
}

func NewVirtualProperty(proxy properties.DeviceProxy) *yeelight.VirtualAdapter {
	return &yeelight.VirtualAdapter{}
}
