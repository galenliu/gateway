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

// HService 不同的service实现此接口
type HService interface {
	GetService() *service.Service
	GetID() string
}

type BridgeProxy struct {
	bridge          *accessory.Bridge
	things          *models.Things
	homekitServices map[string]HService
}

func NewHomeKitBridge(name, sn, manufacturer, model, storagePath string) {

	var bridge *accessory.Bridge
	Bridge = &BridgeProxy{bridge, models.NewThings(), make(map[string]HService)}

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

// GetThings get things form things models,add service to bridge
func (p *BridgeProxy) GetThings() {
	mapOfThings := p.things.GetThings()
	for _, thing := range mapOfThings {
		s := NewHomekitService(thing)
		p.addService(s)
	}
}

// bridge append service
func (p *BridgeProxy) addService(s HService) {
	p.homekitServices[s.GetID()] = s
	p.bridge.AddService(s.GetService())
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

func (p *BridgeProxy) Start() error {
	p.things.Subscribe(util.ThingAdded, Bridge.onThingAdded)
	p.things.Subscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	p.things.Subscribe(util.ThingModified, Bridge.onThingsModified)
	start()
	return nil
}

func (p *BridgeProxy) Stop() {
	p.things.Unsubscribe(util.ThingAdded, Bridge.onThingAdded)
	p.things.Unsubscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	p.things.Unsubscribe(util.ThingModified, Bridge.onThingsModified)
	stop()
}
