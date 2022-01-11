package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

var on = "on"
var level = "level"
var color = "color"

type YeelightProperty struct {
	bulb *yeelight.Yeelight
	*proxy.Property
}

func NewYeelightProperty(bulb *yeelight.Yeelight, description properties.PropertyDescription) *YeelightProperty {
	return &YeelightProperty{bulb, proxy.NewProperty(description)}
}

func (p *YeelightProperty) SetValue(v interface{}) {
	switch p.GetName() {
	case on:
		b, ok := v.(bool)
		if ok {
			err := p.bulb.SetPower(b)
			if err != nil {
				fmt.Printf("turn on error:%s \t\n", err)
				return
			}
		}
		fmt.Printf("set on : %v \t\n", v)
		p.SetCachedValue(v)
		p.NotifyChanged()
		break
	case level:
		f, ok := v.(float64)
		if ok {
			_, err := p.bulb.SetBrightness(int(f))
			if err != nil {
				fmt.Printf("turn level error:%s \t\n", err)
				return
			}
		}
		fmt.Printf("set level : %v \t\n", v)
		p.SetCachedValue(v)
		p.NotifyChanged()
		return

	case color:
		f, ok := v.(string)
		if ok {
			c, err := addon.HTMLToRGB(f)
			if err != nil {
				return
			}
			_, err = p.bulb.SetRGB(c)
			if err != nil {
				fmt.Printf("turn on error:%s \t\n", err)
				return
			}
		}
		fmt.Printf("set color : %v \t\n", v)
		p.SetCachedValue(v)
		p.NotifyChanged()
		return

	}
}
