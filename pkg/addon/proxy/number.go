package proxy

import (
	"fmt"
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
	f, ok := v.(properties.Number)
	if !ok {
		fmt.Printf("value error:%s", v)
		return
	}
	if min := p.GetMinValue(); min > f {
		f = min
	}
	if max := p.GetMinValue(); max < f {
		f = max
	}
	if ok {
		p.NumberInstance.SetValue(f)
	}
}
