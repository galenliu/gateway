package homekit

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/plugin"
	json "github.com/json-iterator/go"
)

var Bridge *BridgeProxy
var stop func()
var start func()

type BridgeProxy struct {
	bridge *accessory.Bridge
}

func NewHomeKitBridge(name, sn, manufacturer, model, storagePath string) {

	var bridge *accessory.Bridge
	Bridge = &BridgeProxy{bridge}

	info := accessory.Info{
		Name:             name,
		SerialNumber:     sn,
		Manufacturer:     manufacturer,
		Model:            model,
		FirmwareRevision: FirmwareRevision,
		ID:               0,
	}
	config := hc.Config{
		StoragePath: storagePath,
		Pin:         "1234432312",
	}
	bridge = accessory.NewBridge(info)

	t, err := hc.NewIPTransport(config, Bridge.bridge.Accessory)
	if err != nil {
		fmt.Printf("create homekit transport err:", err)
		return
	}
	stop = func() {
		<-t.Stop()
	}
	start = func() {
		t.Start()
	}

}

func (p *BridgeProxy) GetServices() {
	devices := plugin.GetDevices()
	for _, dev := range devices {
		_, err := json.MarshalIndent(dev, "", "  ")
		if err != nil {
			log.Error(err.Error())
			continue
		}

	}
}

func (p *BridgeProxy) AddService(s *service.Service) {
	p.bridge.AddService(s)
}

func (p *BridgeProxy) Start() error {
	start()
	return nil
}

func (p *BridgeProxy) Stop() {
	stop()
}
