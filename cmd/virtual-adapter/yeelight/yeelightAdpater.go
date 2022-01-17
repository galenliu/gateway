package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon/proxy"
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

func (a *YeelightAdapter) StartPairing(timeout time.Duration) {
	devices := make(map[string]proxy.DeviceProxy, 1)
	discover := func() {
		bulb, err := yeelight.Discover()
		if err != nil {
			fmt.Printf("adapter:%s StartPairing err:%s \t\n", a.GetId(), err.Error())
			return
		}
		_, ok := devices[bulb.GetAddr()]
		if !ok {
			devices[bulb.GetAddr()] = proxy.NewDevice(NewYeelightBulb(bulb))
		}
	}
	discover()
	for _, d := range devices {
		a.AddDevices(d)
	}

}
