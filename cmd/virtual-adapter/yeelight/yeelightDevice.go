package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/xiam/to"
	"strconv"
)

type YeelightDevice struct {
	*devices.Light
	*yeelight.Yeelight
}

func NewYeelightBulb(bulb *yeelight.Yeelight) *YeelightDevice {
	yeeDevice := &YeelightDevice{
		Light:    devices.NewLightBulb(bulb.GetAddr()), //proxy.NewDevice([]string{schemas.CapabilityLight, schemas.CapabilityOnOffSwitch}, bulb.GetAddr(), "yeelight"+bulb.GetAddr()),
		Yeelight: bulb,
	}
	for _, method := range bulb.GetSupports() {
		switch method {
		case "set_power":
			propValue := bulb.GetPropertyValue("power")
			var value = false
			if propValue != nil {
				v, ok := propValue.(string)
				if ok {
					if v == "on" {
						value = true
					}
				}
			}
			prop := NewOn(bulb, value)
			yeeDevice.AddProperties(prop)
		case "set_bright":
			propValue := bulb.GetPropertyValue("bright")
			var value properties.Integer = 0
			if propValue != nil {
				s, ok := propValue.(string)
				if ok {
					v, err := strconv.Atoi(s)
					if err == nil {
						value = properties.Integer(v)
					}
				}
			}
			prop := NewBrightness(bulb, value)
			yeeDevice.AddProperties(prop)
		case "set_rgb":

			propValue := bulb.GetPropertyValue("rgb")
			value := "#ffffff"
			if propValue != nil {
				_, ok := propValue.(string)
				if ok {
					color := "#" + fmt.Sprintf("%X", to.Int64(propValue))
					value = color
				}
			}
			prop := NewColor(bulb, value)
			yeeDevice.AddProperties(prop)
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
			for n, v := range msg.Params {
				if n == "power" {
					b := v == "on"
					d.GetProperty("on").SetCachedValueAndNotify(b)
				}
				if n == "bright" {
					d.GetProperty("bright").SetCachedValueAndNotify(v)
				}
				if n == "rgb" {
					i := to.Int64(v)
					v := "#" + fmt.Sprintf("%X", i)
					d.GetProperty("color").SetCachedValueAndNotify(v)
				}
			}
		}
	}
}

func (d *YeelightDevice) SetPin(pin string) error {
	fmt.Printf("device: %s set pin: %s \n", d.GetId(), pin)
	return nil
}
