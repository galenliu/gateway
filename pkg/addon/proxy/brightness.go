package proxy

import "github.com/galenliu/gateway/pkg/addon/devices"

type BrightnessInstance interface {
	devices.PropertyEntity
	SetBrightness(int)
}

type BrightnessProperty struct {
	BrightnessInstance
}

func NewBrightness(p BrightnessInstance) *BrightnessProperty {
	return &BrightnessProperty{p}
}

func (p *BrightnessProperty) SetValue(v any) {
	value, ok := v.(int)
	if ok {
		p.SetBrightness(value)
	}
}