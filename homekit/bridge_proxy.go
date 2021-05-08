package homekit

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
)

var Bridge *BridgeProxy
var stop func()
var start func()

type BridgeProxy struct {
	bridge *accessory.Bridge
	Things *models.Things
}

func NewHomeKitBridge(name, sn, manufacturer, model, storagePath string) {

	var bridge *accessory.Bridge
	Bridge = &BridgeProxy{bridge, models.NewThings()}

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

func (p *BridgeProxy) GetThings() {
	p.Things.GetThings()
}

func (p *BridgeProxy) updateServices() {

}

func (p *BridgeProxy) onThingAdded(data []byte) {
	deviceId := json.Get(data, "deviceId").ToString()
	log.Info(deviceId)

}

func (p *BridgeProxy) onThingsModified(data []byte) {

}

func (p *BridgeProxy) onPropertyChanged(data []byte) {

}

func (p *BridgeProxy) AddService(s *service.Service) {
	p.bridge.AddService(s)
}

func (p *BridgeProxy) Start() error {
	p.Things.Subscribe(util.ThingAdded, Bridge.onThingAdded)
	p.Things.Subscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	p.Things.Subscribe(util.MODIFIED, Bridge.onThingsModified)
	start()
	return nil
}

func (p *BridgeProxy) Stop() {
	p.Things.Unsubscribe(util.ThingAdded, Bridge.onThingAdded)
	p.Things.Unsubscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	p.Things.Unsubscribe(util.MODIFIED, Bridge.onThingsModified)
	stop()
}
