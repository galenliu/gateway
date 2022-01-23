package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"github.com/xiam/to"
)

type YeelightDevice struct {
	*devices.LightBulb
	*yeelight.Yeelight
}

func NewYeelightBulb(bulb *yeelight.Yeelight) *YeelightDevice {
	yeeDevice := &YeelightDevice{
		LightBulb: devices.NewLightBulb(devices.DeviceDescription{
			Id: bulb.GetAddr(),
		}), //proxy.NewDevice([]string{schemas.Light, schemas.OnOffSwitch}, bulb.GetAddr(), "yeelight"+bulb.GetAddr()),
		Yeelight: bulb,
	}
	for _, method := range bulb.GetSupports() {
		switch method {
		case "set_power":
			prop := NewOn(bulb)
			yeeDevice.AddProperty(proxy.NewBooleanProxy(prop))
		case "set_bright":
			prop := NewBrightness(bulb)
			yeeDevice.AddProperty(proxy.NewIntegerProxy(prop))
		case "set_rgb":
			prop := NewColor(bulb)
			yeeDevice.AddProperty(proxy.NewStringProxy(prop))
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
					b := v == "on"

					d.GetPropertyEntity("on").SetCachedValueAndNotify(b)
				}
				if n == "bright" {
					d.GetPropertyEntity("bright").SetCachedValueAndNotify(v)
				}
				if n == "rgb" {
					i := to.Int64(v)
					v := "#" + fmt.Sprintf("%X", i)
					d.GetPropertyEntity("color").SetCachedValueAndNotify(v)
				}
			}
		}
	}
}

func (d *YeelightDevice) SetPin(pin string) error {
	fmt.Printf("device: %s set pin: %s \n", d.GetId(), pin)
	return nil
}
