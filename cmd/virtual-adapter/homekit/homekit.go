package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/integration"
	"github.com/galenliu/gateway/pkg/addon/schemas"
)

type manager interface {
	send()
}

type HomeKit struct {
	*integration.Integration
	*accessory.Bridge
	manager     manager
	storagePath string
}

func NewHomekitIntegration() *HomeKit {
	info := accessory.Info{
		Name:             "Wot Gateway",
		SerialNumber:     "12344321",
		Manufacturer:     "WOT",
		Model:            "WebThings",
		FirmwareRevision: "1.1",
		ID:               0,
	}
	hk := &HomeKit{
		Integration: nil,
		Bridge:      accessory.NewBridge(info),
		manager:     nil,
	}
	hk.start()
	return hk
}

func (c *HomeKit) GetThings() {
	ts := c.Integration.LoadThings()
	for _, t := range ts {
		thing := NewThing(t)
		c.AddThing(thing)
	}
}

func (c *HomeKit) OnThing(t things.Thing) {
	for _, s := range t.AtType {
		switch s {
		case schemas.Light:

		}
	}
}

func (c *HomeKit) register() {

}

func (c *HomeKit) start() {
	t, err := hc.NewIPTransport(hc.Config{
		StoragePath: c.storagePath,
		Port:        "12345",
		Pin:         "12344321",
		SetupId:     "",
	}, c.Bridge.Accessory)

	t.Start()
	if err != nil {
		return
	}
	t.Start()
}
