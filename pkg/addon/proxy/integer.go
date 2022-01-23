package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type IntegerInstance interface {
	properties.IntegerEntity
	SetValue(v properties.Integer)
}

type IntegerProxy struct {
	IntegerInstance
}

func NewIntegerProxy(p IntegerInstance) *IntegerProxy {
	return &IntegerProxy{p}
}

func (p *IntegerProxy) SetValue(v any) {
	f, ok := v.(float64)
	value := properties.Integer(f)
	if ok {
		if min := p.GetMinValue(); min > value {
			value = min
		}
		if max := p.GetMaxValue(); max < value {
			value = max
		}
	}
	p.IntegerInstance.SetValue(value)
}
