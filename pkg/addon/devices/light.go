package devices

import (
	"strconv"
	"strings"
)

type LightHandler interface {
	TurnOn()
	TurnOff()
	SetBrightness(brightness int)
}

type Light struct {
	*Device
}

func NewLightBulb(id string, args ...string) *Light {
	if id == "" {
		return nil
	}
	title := "light" + id
	desc := ""
	if len(args) > 0 {
		title = args[0]
	}
	if len(args) > 1 {
		desc = args[1]
	}
	return &Light{
		NewDevice(DeviceDescription{
			Id:          id,
			AtType:      []Capability{CapabilityLight, CapabilityOnOffSwitch},
			Title:       title,
			Description: desc,
		}),
	}
}

func (light *Light) TurnOn() {
	//light.On.SetValue(true)
}

func (light *Light) TurnOff() {
	//light.On.SetValue(false)
}

func (light *Light) Toggle() {
	//if light.On.Value == true {
	//	light.TurnOff()
	//} else {
	//	light.TurnOn()
	//}
}

func (light *Light) SetBrightness(brightness int) {
	//if light.Bright == nil {
	//	return
	//}
	//if brightness == 0 && light.On.Value == true {
	//	light.TurnOff()
	//} else if brightness > 0 && light.On.Value == false {
	//	light.TurnOn()
	//}
	//light.Bright.SetValue(brightness)
}

func (light *Light) propertyValueUpdate(propName string, newValue any) {

}

func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}
