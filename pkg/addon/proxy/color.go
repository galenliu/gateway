package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
	"image/color"
)

type ColorPropertyInstance interface {
	properties.Entity
	SetValue(c color.RGBA)
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
		if err == nil {
			p.ColorPropertyInstance.SetValue(c)
		}
	}
}
