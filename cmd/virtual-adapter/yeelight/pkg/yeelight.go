package yeelight

import (
	"fmt"
	c "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg/color"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg/utils"
	"image/color"
	"log"
	"net"
)

type EffectType string

const (
	Smooth EffectType = "smooth"
	Sudden            = "sudden"
)

type LightType int

const (
	Main    LightType = 0
	Ambient           = 1
)

type Mode int

const (
	Last      Mode = 0
	Normal         = 1
	RGB            = 2
	HSV            = 3
	ColorFlow      = 4
	Moonlight      = 5
)

type (
	PropsResult struct {
		ID     int
		Result map[string]string
		Error  *Error
	}
)

//Bulb represents device
type Bulb struct {
	conn   net.Conn
	params *YeelightParams
	ip     string
	addr   string
	effect EffectType
	cmdId  int
}

func (y *Yeelight) TurnOn() (*CommandResult, error) {
	return y.executeCommand("set_power", "on")
}

func (y *Yeelight) TurnOnWithParams(mode Mode, duration int) (*CommandResult, error) {
	return y.executeCommand("set_power", "on", duration, mode)
}

func (y *Yeelight) TurnOff() (*CommandResult, error) {
	return y.executeCommand("set_power", "off")
}

func (y *Yeelight) IsOn() (bool, error) {
	res, err := y.GetProps([]string{"power"})
	if err != nil {
		return false, err
	}
	power := res.Result["power"]

	return power == "on", nil
}

func (y *Yeelight) SetBrightness(brightness int) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			y.SetPower(true)
		}
	}
	return y.executeCommand("set_bright", utils.GetBrightnessValue(brightness))
}

func (y *Yeelight) SetRGB(rgba color.RGBA) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			y.SetPower(true)
		}
	}
	value := c.RGBToYeelight(rgba)
	return y.executeCommand("set_rgb", value)
}

func (y *Yeelight) SetHSV(hue int, saturation int) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			y.SetPower(true)
		}
	}
	return y.executeCommand("set_rgb", hue, saturation)
}

func (y *Yeelight) SetBrightnessWithDuration(brightness int, duration int) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			y.SetPower(true)
		}
	}
	if !checkBrightnessValue(brightness) {
		log.Fatalln("The brightness value to set (1-100)")
	}

	return y.executeCommand("set_bright", brightness, duration)
}

func (y *Yeelight) StartFlow(flow *Flow) (*CommandResult, error) {

	params := flow.AsStartParams()
	return y.executeCommand("start_cf", params)
}

func (y *Yeelight) StopFlow() (*CommandResult, error) {
	return y.executeCommand("stop_cf", "")
}

func (y *Yeelight) GetProps(props []string) (*PropsResult, error) {
	res, err := y.executeCommand("get_prop", props)
	if err != nil {
		return nil, err
	}

	propsMap := make(map[string]string)

	for i, val := range res.Result {
		key := props[i]
		propsMap[key] = fmt.Sprintf("%v", val)
	}

	return &PropsResult{ID: res.ID, Error: res.Error, Result: propsMap}, nil
}

func (y *Yeelight) GetAddr() string {
	return y.addr
}

func (y *Yeelight) SetName(name string) (*CommandResult, error) {
	return y.executeCommand("set_name", name)
}
