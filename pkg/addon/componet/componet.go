package componet

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Component struct {
	*accessory.Bridge
}

func (c *Component) AddThing(s *service.Service) {
	c.Bridge.AddService(s)
}

func (c *Component) Run() {
	config := hc.Config{
		StoragePath: "",
		Port:        "",
		Pin:         "",
		SetupId:     "",
	}
	transport, err := hc.NewIPTransport(config, c.Accessory)
	if err != nil {
		return
	}
	transport.Start()
}
