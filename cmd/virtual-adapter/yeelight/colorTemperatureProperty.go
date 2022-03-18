package yeelight

import (
	"context"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/lib"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"time"
)

type ColorTemperatureProperty struct {
	device *YeelightDevice
	*properties.ColorTemperatureProperty
}

func NewColorTemperatureProperty(bulb *YeelightDevice, value properties.Integer) *ColorTemperatureProperty {
	return &ColorTemperatureProperty{
		device:                   bulb,
		ColorTemperatureProperty: properties.NewColorTemperatureProperty(value),
	}
}

func (ct ColorTemperatureProperty) SetValue(v properties.Integer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err := ct.device.Client.SetColorTemperature(ctx, int(v), yeelight.EffectSmooth, duration)
	if err != nil {
		return err
	}
	ct.SetCachedValue(v)
	ct.NotifyChanged()
	return nil
}
