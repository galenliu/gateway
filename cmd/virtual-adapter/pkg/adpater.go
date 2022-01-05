package pkg

import (
	"fmt"
	"github.com/akominch/yeelight"
	"github.com/galenliu/gateway/pkg/addon"
	"time"
)

type VirtualAdapter struct {
	*addon.Adapter
}

func NewVirtualAdapter(manager *addon.Manager, adapterId, name string) *VirtualAdapter {
	v := &VirtualAdapter{
		addon.NewAdapter(manager, adapterId, name),
	}
	v.StartPairing(300 * time.Duration(time.Millisecond))
	return v
}

func (a *VirtualAdapter) StartPairing(timeout time.Duration) {
	bulb, err := yeelight.Discover()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("bulb: %v", bulb)
}
