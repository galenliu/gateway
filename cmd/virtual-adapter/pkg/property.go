package pkg

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type VirtualProperty struct {
	*properties.Property
}

func NewVirtualProperty(proxy properties.DeviceProxy) *VirtualAdapter {
	return &VirtualAdapter{}
}
