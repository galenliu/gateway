package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Property struct {
	*properties.Property
}

func NewProperty(des properties.PropertyDescription) *Property {
	return &Property{
		Property: properties.NewProperty(des),
	}
}

func (p *Property) SetHandler(handler addon.DeviceProxy) {
	p.Property.Handler = handler
}

func (p *Property) SetValue(v interface{}) {
	fmt.Print("property: %s SetValue func not implemented", p.GetName())
}
