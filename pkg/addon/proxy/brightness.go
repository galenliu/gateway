package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type BrightnessInstance interface {
	properties.Entity
	SetBrightness(int)
}

type BrightnessProperty struct {
	BrightnessInstance
}

func NewBrightness(p BrightnessInstance) *BrightnessProperty {
	return &BrightnessProperty{p}
}

func (p *BrightnessProperty) SetValue(v any) {
	value, ok := v.(float64)
	if ok {
		p.SetBrightness(int(value))
	}
}
