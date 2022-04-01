package yeelight

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/lib"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
	"log"
	"strconv"
	"strings"
	"time"
)

const timeOut = 1000 * time.Millisecond
const duration = 100 * time.Millisecond

type YeelightDevice struct {
	*devices.Light
	ctx context.Context
	*yeelight.Client
	id   string
	name string
	ip   string
}

func NewYeelightBulb(clint *yeelight.Client, id, name, ip string) *YeelightDevice {
	yeeDevice := &YeelightDevice{
		id:     id,
		ctx:    context.Background(),
		name:   name,
		ip:     ip,
		Light:  devices.NewLightBulb(id), //proxy.NewDevice([]string{schemas.CapabilityLight, schemas.CapabilityOnOffSwitch}, bulb.GetAddr(), "yeelight"+bulb.GetAddr()),
		Client: clint,
	}

	props, err := yeeDevice.Client.GetProperties(context.Background(), []string{yeelight.PropertyPower, yeelight.PropertyColorMode, yeelight.PropertyBright, yeelight.PropertyCT, yeelight.PropertyHue, yeelight.PropertyRGB})
	if err != nil {
		log.Fatalln(err)
	}
	for name, value := range props {
		fmt.Println("> ", name, ":", value)
	}
	for name, value := range props {
		if value == "" {
			continue
		}
		switch name {
		case yeelight.PropertyPower:
			prop := NewOn(yeeDevice, isPowerOn(value))
			yeeDevice.AddProperties(prop)
		case yeelight.PropertyBright:

			value, _ := strconv.Atoi(value)
			prop := NewBrightness(yeeDevice, properties.Integer(value))
			yeeDevice.AddProperties(prop)
		case yeelight.PropertyRGB:
			prop := NewColor(yeeDevice, value)
			yeeDevice.AddProperties(prop)
		case yeelight.PropertyColorMode:
			prop := NewColorMode(yeeDevice, value)
			yeeDevice.AddProperties(prop)
		case yeelight.PropertyCT:
			v, err := strconv.Atoi(value)
			if err != nil {
				break
			}
			prop := NewColorTemperatureProperty(yeeDevice, properties.Integer(v))
			yeeDevice.AddProperties(prop)
		default:
			continue
		}
	}
	go yeeDevice.Listen(yeeDevice.ctx)
	return yeeDevice
}

func (d *YeelightDevice) SetCredentials(username, password string) error {
	return nil
}

type NotifyMessage struct {
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
}

func (d *YeelightDevice) Listen(ctx context.Context) error {
	channel, err := d.Client.Listen(d.ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		select {
		case data := <-channel:
			var notify NotifyMessage
			data = strings.Trim(data, "\r\n")
			err := json.Unmarshal([]byte(data), &notify)
			fmt.Printf("yeelight device notify:%s \t\n", data)
			if err == nil {
				if notify.Method == "props" {
					for name, value := range notify.Params {
						switch name {
						case yeelight.PropertyPower:
							d.OnOff.SetCachedValueAndNotify(isPowerOn(to.String(value)))
						case yeelight.PropertyBright:
							v := to.Float64(value)
							if err == nil {
								d.Brightness.SetCachedValueAndNotify(v)
							}
						case yeelight.PropertyColorMode:
							if value == "2" {
								d.ColorMode.SetCachedValueAndNotify(properties.ColorModePropertyEnumColor)
							}
							if value == "3" {
								d.ColorMode.SetCachedValueAndNotify(properties.ColorModePropertyEnumTemperature)
							}
						case yeelight.PropertyRGB:
							v := to.Float64(value)
							s := strconv.FormatUint(uint64(v), 16)
							d.Color.SetCachedValueAndNotify("#" + s)
						case yeelight.PropertyCT:
							v := to.Float64(value)
							if err != nil {
								d.ColorTemperature.SetCachedValueAndNotify(v)
							}
						case yeelight.BrightWithZero:
							v := to.Float64(value)
							if err != nil {
								d.Brightness.SetCachedValueAndNotify(v)
							}
						default:
							fmt.Printf("Bad Params name: %s,value : %v \t\n", name, value)
						}
					}
				}
			}
		case <-ctx.Done():
			fmt.Printf("exit")
		}
	}
}

func (d *YeelightDevice) SetPin(pin string) error {
	fmt.Printf("device: %s set pin: %s \n", d.GetId(), pin)
	return nil
}

func (d *YeelightDevice) HandleRemoved() {
	d.ctx.Done()
}
