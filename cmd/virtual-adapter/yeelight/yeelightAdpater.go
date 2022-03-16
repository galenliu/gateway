package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type YeelightAdapter struct {
	*proxy.Adapter
}

func NewVirtualAdapter(adapterId string) *YeelightAdapter {
	v := &YeelightAdapter{
		proxy.NewAdapter(adapterId, "yeelight"),
	}
	return v
}

func (a *YeelightAdapter) HandleDeviceSaved(msg messages.DeviceSavedNotificationJsonData) {
	fmt.Printf("yeelight-adapter handle device saved deviceId: %s \t\n", msg.DeviceId)
}

func (a *YeelightAdapter) StartPairing(timeout <-chan time.Time) {
	fmt.Printf("start pairing...\n")
	if timeout == nil {
		timeout = time.After(5 * time.Second)
	}

	discover := func() *yeelight.Yeelight {
		bulb, err := yeelight.Discover()
		if err != nil {
			fmt.Printf("adapter:%s  err:%s \t\n", a.GetId(), err.Error())
			return nil
		}
		return bulb
	}
	for {
		foundDevice := discover()
		if foundDevice != nil {
			yl := NewYeelightBulb(foundDevice)
			a.HandleDeviceAdded(yl)
		}
		select {
		case <-timeout:
			return
		}
	}

}
