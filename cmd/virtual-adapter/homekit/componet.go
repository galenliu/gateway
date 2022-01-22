package homekit

import (
	"github.com/brutella/hc/accessory"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/integration"
)

type manager interface {
	send()
}

type Component struct {
	*integration.Integration
	*accessory.Bridge
	manager manager
}

func (c *Component) OnThings(things []things.Thing) {

}

func (c *Component) send(msg any) {

}

func (c *Component) register() {

}

func (c *Component) start() {

}
