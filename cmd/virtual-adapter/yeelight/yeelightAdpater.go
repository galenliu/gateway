package yeelight

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"time"
)

type VirtualAdapter struct {
	*proxy.Adapter
}

func NewVirtualAdapter(adapterId, packageName string) *VirtualAdapter {
	manager, err := proxy.NewAddonManager(packageName)
	if err != nil {
		fmt.Printf("addon manager error: %s", err.Error())
		return nil
	}
	v := &VirtualAdapter{
		proxy.NewAdapter(manager, adapterId, packageName),
	}
	manager.AddAdapters(v)
	v.StartPairing(300 * time.Duration(time.Millisecond))
	return v
}

func (a *VirtualAdapter) StartPairing(timeout time.Duration) {

	devices := make(map[string]addon.DeviceProxy, 1)

	discover := func() {
		bulb, err := yeelight.Discover()
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		_, ok := devices[bulb.GetAddr()]
		if !ok {
			devices[bulb.GetAddr()] = NewYeelightBulb(bulb)
		}
	}
	discover()
	for _, d := range devices {
		a.AddDevices(d)
	}

}
