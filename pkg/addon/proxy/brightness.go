package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type BrightnessInstance interface {
	properties.Entity
	SetBrightness(int) error
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
		err := p.SetBrightness(int(value))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		p.SetCachedValueAndNotify(value)
	}
}
