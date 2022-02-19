package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type StringInstance interface {
	properties.StringEntity
	SetValue(v string) error
}

type StringProxy struct {
	StringInstance
}

func NewStringProxy(p StringInstance) *StringProxy {
	return &StringProxy{p}
}

func (p *StringProxy) SetValue(v any) {
	value := p.CheckValue(v)
	err := p.StringInstance.SetValue(value)
	if err != nil {
		fmt.Printf("device:%s set property:%s value error: %v", p.GetDevice().GetId(), p.GetName(), err)
	}
}
