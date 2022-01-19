package homekit

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/schemas"

	"github.com/brutella/hc/service"
)

type Thing struct {
	*things.Thing

}

func (t *Thing) ToAccessory() {
	if t.AtType
		service.New()
}

func GetService(t string) *service.Service {
	switch t {
	case schemas.Light:
		light := service.NewLightbulb()
		light.Service
		return light.Service
	}
	return nil
}
