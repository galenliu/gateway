package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/xiam/to"
	"image/color"
)

type ColorPropertyInstance interface {
	properties.Entity
	SetValue(c color.RGBA) error
}

type ColorProperty struct {
	ColorPropertyInstance
}

func NewColor(p ColorPropertyInstance) *ColorProperty {
	return &ColorProperty{p}
}

func (p *ColorProperty) SetValue(v any) {
	value, ok := v.(string)
	if ok {
		c, err := HTMLToRGB(value)
		if err != nil {
			err := p.ColorPropertyInstance.SetValue(c)
			if err != nil {
				p.SetCachedValueAndNotify(to.String(v))
				return
			}
		}
	}
}
