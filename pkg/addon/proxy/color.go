package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type ColorPropertyInstance interface {
	properties.Entity
	SetValue(string2 string)
}

type ColorProperty struct {
	ColorPropertyInstance
}

func NewColor(p ColorPropertyInstance) *StringProperty {
	return &StringProperty{p}
}

func (p *ColorProperty) SetValue(v any) {
	value, ok := v.(string)
	if ok {
		p.ColorPropertyInstance.SetValue(value)
	}
}
