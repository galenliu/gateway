package things

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/homekit"
)

type LightBulb struct {
	*homekit.Thing
	On characteristic.On
}

func NewLightBulb(data []byte) *LightBulb {
	_ = service.NewLightbulb()
	return nil
}
