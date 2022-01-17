package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type StringPropertyInstance interface {
	properties.Entity
	SetValue(string2 string)
}

type StringProperty struct {
	StringPropertyInstance
}

func NewString(p StringPropertyInstance) *StringProperty {
	return &StringProperty{p}
}

func (p *StringProperty) SetValue(v any) {
	value, ok := v.(string)
	if ok {
		p.StringPropertyInstance.SetValue(value)
	}
}
