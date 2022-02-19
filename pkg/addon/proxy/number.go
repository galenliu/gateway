package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type NumberInstance interface {
	properties.NumberEntity
	SetValue(v properties.Number)
}

type NumberProxy struct {
	NumberInstance
}

func NewNumberProxy(p NumberInstance) *NumberProxy {
	return &NumberProxy{p}
}

func (p *NumberProxy) SetValue(v any) {
	value := p.CheckValue(v)
	p.NumberInstance.SetValue(value)

}
