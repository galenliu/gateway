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

	discover := func(done <-chan time.Time) {
		for {
			select {
			case <-done:
				return
			default:
				bulb, err := yeelight.Discover()
				if err != nil {
					fmt.Printf(err.Error())
					continue
				}
				p, err := bulb.Discover()
				if err != nil {
					fmt.Printf(err.Error())
					continue
				}
				lihgt := New
				fmt.Printf("bulb: %v", bulb)
			}
		}
	}
	done := time.After(timeout)
	go discover(done)
}
