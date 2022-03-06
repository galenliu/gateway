package virtual

import "github.com/galenliu/gateway/pkg/addon/properties"

type Property struct {
	*properties.Property
}

func NewProperty(p properties.Entity) *Property {
	return &Property{p.GetProperty()}
}

func (p *Property) SetValue(v properties.Number) error {
	return nil

}
