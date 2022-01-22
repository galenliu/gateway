package service

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/schemas"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
)

type LightBulb struct {
	light *service.Lightbulb
	thing *things.Thing
}

func NewLightBulb(t *things.Thing) *LightBulb {
	l := &LightBulb{thing: t}
	p := l.getProperty(schemas.OnOffProperty)
	if p == nil {
		return nil
	}
	l.light.On = characteristic.NewOn()

	return l
}

func (b *LightBulb) getProperty(typ string) wot.PropertyAffordance {
	for _, p := range b.thing.Properties {
		if p.GetAtType() == typ {
			return p
		}
	}
	return nil
}
