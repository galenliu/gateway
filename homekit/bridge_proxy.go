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
	*accessory.Bridge
	things          *models.Things
	homekitServices map[string]HService
	info            accessory.Info
	config          hc.Config
}

func NewHomeKitBridge(name, sn, manufacturer, model, storagePath, pin string) {

	b := &BridgeProxy{}
	b.info = accessory.Info{
		Name:             name,
		SerialNumber:     sn,
		Manufacturer:     manufacturer,
		Model:            model,
		FirmwareRevision: FirmwareRevision,
		ID:               0,
	}
	b.config = hc.Config{
		StoragePath: storagePath,
		Pin:         pin,
	}
	b.Bridge = accessory.NewBridge(b.info)
}

// GetThings get things form things models,add service to bridge
func (br *BridgeProxy) GetThings() {
	mapOfThings := br.things.GetThings()
	for _, thing := range mapOfThings {
		s := NewHomekitService(thing)
		br.addService(s)
	}
}

// bridge append service
func (br *BridgeProxy) addService(s HService) {
	br.homekitServices[s.GetID()] = s
	br.Bridge.AddService(s.GetService())
}

func (br *BridgeProxy) updateServices() {
	br.Services = make([]*service.Service, 1)
	br.GetThings()
}

func (br *BridgeProxy) onThingAdded(data []byte) {
	deviceId := json.Get(data, "deviceId").ToString()
	log.Info(deviceId)

}

func (br *BridgeProxy) onThingsModified(data []byte) {

}

func (br *BridgeProxy) onPropertyChanged(data []byte) {

}

func (br *BridgeProxy) Start() error {

	t, err := hc.NewIPTransport(br.config, br.Bridge.Accessory)

	if err != nil {
		return fmt.Errorf("create homekit transport err:", err)

	}
	stop = func() {
		<-t.Stop()
	}
	t.Start()
	br.things.Subscribe(util.ThingAdded, Bridge.onThingAdded)
	br.things.Subscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	br.things.Subscribe(util.ThingModified, Bridge.onThingsModified)
	start()
	return nil
}

func (br *BridgeProxy) Stop() {
	br.things.Unsubscribe(util.ThingAdded, Bridge.onThingAdded)
	br.things.Unsubscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	br.things.Unsubscribe(util.ThingModified, Bridge.onThingsModified)
	stop()
}
