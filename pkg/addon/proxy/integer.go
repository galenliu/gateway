package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type IntegerInstance interface {
	properties.IntegerEntity
	SetValue(v properties.Integer) error
}

type IntegerProxy struct {
	IntegerInstance
}

func NewIntegerProxy(p IntegerInstance) *IntegerProxy {
	return &IntegerProxy{p}
}

func (p *IntegerProxy) SetValue(v any) {
	value := p.CheckValue(v)
	err := p.IntegerInstance.SetValue(value)
	if err != nil {
		fmt.Printf("device:%s set property:%s value error: %v", p.GetDevice().GetId(), p.GetName(), err)
	}
}
