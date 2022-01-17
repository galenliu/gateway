package proxy

import (
	"fmt"
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

func (p *Property) SetHandler(handler DeviceProxy) {
	p.Property.SetHandler(handler)
}

func (p *Property) SetValue(v any) {
	fmt.Printf("property: %s SetValue func not implemented", p.GetName())
}
