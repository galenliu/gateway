package homekit

import (
	"gateway/pkg/log"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

var HCBridge *BridgeProxy

type BridgeProxy struct {
	bridge *accessory.Bridge
}

func NewHomeKitBridge(name, sn, manufacturer, model, storagePath string) {

	var bridge *accessory.Bridge
	HCBridge = &BridgeProxy{bridge}

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

	t, err := hc.NewIPTransport(config, HCBridge.bridge.Accessory)
	if err != nil {
		log.Error("create homekit transport err:", err)
		return
	}
	hc.OnTermination(func() {
		<-t.Stop()
	})
	t.Start()
}
