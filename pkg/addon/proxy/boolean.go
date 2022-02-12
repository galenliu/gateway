package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type BooleanInstance interface {
	properties.Entity
	TurnOff() error
	TurnOn() error
}

type BooleanProxy struct {
	BooleanInstance
}

func NewBooleanProxy(p BooleanInstance) *BooleanProxy {
	return &BooleanProxy{p}
}

func (p *BooleanProxy) SetValue(v any) {
	value, ok := v.(bool)
	if !ok {
		fmt.Printf("value error:%s", v)
		return
	}
	if value {
		err := p.BooleanInstance.TurnOn()
		if err != nil {
			return
		}
	} else {
		err := p.BooleanInstance.TurnOff()
		if err != nil {
			return
		}
	}
}
