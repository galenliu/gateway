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

type LightBulb struct {
	*Device
	//On     addon.IProperty
	//Bright addon.IProperty
	//Color  addon.IProperty
}

//func NewLightBulb(id, title string) *LightBulb {
//
//	lightBulb := &LightBulb{}
//	lightBulb.Addon_Device = addon.NewDevice(id, title)
//
//	lightBulb.AddTypes(addon.Light, addon.OnOffSwitch)
//	return lightBulb
//}
//
//func (light *LightBulb) addProperty(p addon.IProperty) {
//	if p.GetAtType() == properties.TypeOnOffProperty {
//		light.On = p
//		light.Addon_Device.addProperty(p)
//		return
//	}
//	if p.GetAtType() == properties.TypeBrightnessProperty {
//		light.Bright = p
//		light.Addon_Device.addProperty(p)
//		return
//	}
//	if p.GetAtType() == properties.TypeColorProperty {
//
//		light.Color = p
//		light.Addon_Device.addProperty(p)
//		return
//	}
//	light.Addon_Device.addProperty(p)
//}

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
