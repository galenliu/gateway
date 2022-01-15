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

func (on *Color) SetValue(v string) {

}
