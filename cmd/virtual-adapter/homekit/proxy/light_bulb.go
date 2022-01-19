package proxy

import "github.com/brutella/hc/service"

type LightBulb struct {
	*service.Lightbulb
}

func (light *LightBulb) TurnOn() {
	light.On.SetValue(true)
}

func (light *LightBulb) TurnOff() {
	light.On.SetValue(false)
}
