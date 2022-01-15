package yeelight

import (
	"fmt"
	yeelight "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Brightness struct {
	bulb *yeelight.Yeelight
	*properties.BrightnessProperty
}

func NewBrightness(bulb *yeelight.Yeelight) *Brightness {
	return &Brightness{
		bulb:               bulb,
		BrightnessProperty: properties.NewBrightnessProperty(properties.PropertyDescription{}),
	}
}

func (b Brightness) SetBrightness(v int) {
	_, err := b.bulb.SetBrightness(v)
	if err != nil {
		fmt.Printf("Error setting brightness:%s", err.Error())
	}
}