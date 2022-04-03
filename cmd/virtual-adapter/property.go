package virtual_adapter

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Property struct {
	*properties.Property
}

func NewVirtualProperty(p *properties.Property) *Property {
	return &Property{p}
}

func (p *Property) SetPropertyValue(v any) error {
	p.SetCachedValue(v)
	p.NotifyChanged()
	fmt.Printf("device: %s set property: %s value: %v \t\n", p.GetDevice().GetId(), p.GetName(), v)
	return nil
}
