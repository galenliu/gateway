package yeelight

import (
	yeelight "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type On struct {
	bulb *yeelight.Yeelight
	*properties.OnOffProperty
}

func NewOn(bulb *yeelight.Yeelight, value bool) *On {
	return &On{
		bulb,
		properties.NewOnOffProperty(value),
	}
}

func (on *On) TurnOff() error {
	_, err := on.bulb.TurnOff()
	if err != nil {
		return err
	}
	on.SetCachedValue(false)
	on.NotifyChanged()
	return nil
}

func (on *On) TurnOn() error {
	_, err := on.bulb.TurnOn()
	if err != nil {
		return err
	}
	on.SetCachedValue(true)
	on.NotifyChanged()
	return nil
}
