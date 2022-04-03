package yeelight_adapter

import (
	"context"
	"github.com/galenliu/gateway/cmd/yeelight-adapter/lib"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type On struct {
	device *YeelightDevice
	*properties.OnOffProperty
}

func NewOn(bulb *YeelightDevice, value bool) *On {
	return &On{
		bulb,
		properties.NewOnOffProperty(value),
	}
}

func (on *On) TurnOff() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	err := on.device.Client.Power(ctx, false, yeelight.PowerModeDefault, yeelight.EffectSmooth, duration)
	if err != nil {
		return err
	}
	on.SetCachedValue(false)
	on.NotifyChanged()
	return nil
}

func (on *On) TurnOn() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	err := on.device.Client.Power(ctx, true, yeelight.PowerModeDefault, yeelight.EffectSmooth, duration)
	if err != nil {
		return err
	}
	on.SetCachedValue(true)
	on.NotifyChanged()
	return nil
}
