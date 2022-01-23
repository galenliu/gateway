package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type StringInstance interface {
	properties.Entity
	SetValue(v string)
}

type StringProxy struct {
	StringInstance
}

func NewStringProxy(p StringInstance) *StringProxy {
	return &StringProxy{p}
}

func (p *StringProxy) SetValue(v any) {
	value, ok := v.(string)
	if !ok {
		fmt.Printf("value error:%s", v)
		return
	}
	if ok {
		p.StringInstance.SetValue(value)
	}
}
