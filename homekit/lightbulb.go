package homekit

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway-addon/properties"
)

type LightBulb struct {
	*service.Lightbulb
	*thing
}

func NewLightBulb(s *thing) *LightBulb {
	l := &LightBulb{}
	l.Lightbulb = service.NewLightbulb()
	l.thing.Service = l.Lightbulb.Service
	for _, p := range l.Properties {
		switch p.AtType {
		case properties.TypeOnOffProperty:
			l.Lightbulb.On = characteristic.NewOn()

		}
	}
	return l
}

func (l *LightBulb) FindProperty() {

}
