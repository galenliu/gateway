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
			if d := a.GetDevice(foundDevice.GetAddr()); d == nil {
				yl := NewYeelightBulb(foundDevice)
				a.AddDevices(yl)
			}
		}
		select {
		case <-timeout:
			return
		}
	}

}
