package homekit

import (
	"github.com/brutella/hc/service"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/integration"
)

type Thing struct {
	*integration.ThingEntity
	*service.Service
}

func NewThing(thing things.Thing) *Thing {
	return &Thing{
		integration.NewThingEntity(&thing),
		nil,
	}
}
