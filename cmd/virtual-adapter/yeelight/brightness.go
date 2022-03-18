package yeelight

import (
	"context"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/lib"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"time"
)

type Brightness struct {
	device *YeelightDevice
	*properties.BrightnessProperty
}

func NewBrightness(bulb *YeelightDevice, value properties.Integer) *Brightness {
	return &Brightness{
		device:             bulb,
		BrightnessProperty: properties.NewBrightnessProperty(value),
	}
}

func (b *Brightness) SetValue(v properties.Integer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err := b.device.Client.SetBright(ctx, int(v), yeelight.EffectSmooth, duration)
	if err != nil {
		return err
	}
	b.SetCachedValue(v)
	b.NotifyChanged()
	return nil
}
