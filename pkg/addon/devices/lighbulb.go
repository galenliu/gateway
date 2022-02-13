package devices

import (
	"github.com/galenliu/gateway/pkg/addon/schemas"
	"strconv"
	"strings"
)

type LightHandler interface {
	TurnOn()
	TurnOff()
	SetBrightness(brightness int)
}

type LightBulb struct {
	*Device
}

func NewLightBulb(description DeviceDescription) *LightBulb {
	if description.Id == "" {
		return nil
	}
	if description.Title == "" {
		description.Title = description.Id
	}
	if description.AtType == nil {
		description.AtType = make([]string, 0)
		description.AtType = append(description.AtType, schemas.Light, schemas.OnOffSwitch)
	} else {
		description.AtType = append(description.AtType, schemas.Light, schemas.OnOffSwitch)
	}
	return &LightBulb{
		NewDevice(description),
	}
}

func (light *LightBulb) TurnOn() {
	//light.On.SetValue(true)
}

func (light *LightBulb) TurnOff() {
	//light.On.SetValue(false)
}

func (light *LightBulb) Toggle() {
	//if light.On.Value == true {
	//	light.TurnOff()
	//} else {
	//	light.TurnOn()
	//}
}

func (light *LightBulb) SetBrightness(brightness int) {
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

func (light *LightBulb) propertyValueUpdate(propName string, newValue any) {

}

func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}
