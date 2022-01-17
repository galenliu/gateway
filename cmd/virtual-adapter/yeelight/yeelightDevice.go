package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"github.com/galenliu/gateway/pkg/addon/schemas"
)

type YeelightDevice struct {
	*proxy.Device
	*yeelight.Yeelight
}

func NewYeelightBulb(bulb *yeelight.Yeelight) *YeelightDevice {
	yeeDevice := &YeelightDevice{
		Device:   proxy.NewDevice([]string{schemas.Light, schemas.OnOffSwitch}, bulb.GetAddr(), "yeelight"+bulb.GetAddr()),
		Yeelight: bulb,
	}
	for _, method := range bulb.GetSupports() {
		switch method {
		case "set_power":
			prop := NewOn(bulb)
			yeeDevice.AddProperty(proxy.NewOnOff(prop))
		case "set_bright":
			prop := NewBrightness(bulb)
			yeeDevice.AddProperty(proxy.NewBrightness(prop))
		case "set_rgb":
			prop := NewColor(bulb)
			yeeDevice.AddProperty(proxy.NewColor(prop))
		default:
			continue
		}
	}
	go func() {
		err := yeeDevice.Listen()
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}
	}()
	return yeeDevice
}

func (d *YeelightDevice) SetCredentials(username, password string) error {
	return nil
}

func (d *YeelightDevice) Listen() error {
	notify, _, err := d.Yeelight.Listen()
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-notify:
			fmt.Printf("notify: %s", msg)
			for n, v := range msg.Params {
				if n == "power" {
					d.GetProperty("on").SetCachedValueAndNotify(v)
				}
				if n == "bright" {
					d.GetProperty("bright").SetCachedValueAndNotify(v)
				}
				if n == "rgb" {
					d.GetProperty("color").SetCachedValueAndNotify(v)
				}
			}
		}
	}
	return nil
}
