package yeelight

import (
	yeelight "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Brightness struct {
	bulb *yeelight.Yeelight
	*properties.BrightnessProperty
}

func NewBrightness(bulb *yeelight.Yeelight, value properties.Integer) *Brightness {
	return &Brightness{
		bulb:               bulb,
		BrightnessProperty: properties.NewBrightnessProperty(value),
	}
}

func (b Brightness) SetValue(v properties.Integer) error {
	_, err := b.bulb.SetBrightness(int(v))
	if err != nil {
		return err
	}
	b.SetCachedValue(v)
	b.NotifyChanged()
	return nil
}
