package yeelight

import (
	yeelight "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Color struct {
	bulb *yeelight.Yeelight
	*properties.ColorProperty
}

func NewColor(bulb *yeelight.Yeelight) *Color {
	return &Color{
		bulb,
		properties.NewColorProperty(properties.PropertyDescription{}),
	}
}

func (on *Color) SetValue(v string) error {
	c, err := properties.HTMLToRGB(v)
	if err != nil {
		return err
	}
	_, err = on.bulb.SetRGB(c)
	if err != nil {

		return err
	}
	on.SetCachedValue(v)
	on.NotifyChanged()
	return nil
}
