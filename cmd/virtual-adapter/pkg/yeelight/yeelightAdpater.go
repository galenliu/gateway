package yeelight

import (
	"fmt"
	"github.com/akominch/yeelight"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	uuid "github.com/satori/go.uuid"
	"time"
)

type VirtualAdapter struct {
	*proxy.Adapter
}

func NewVirtualAdapter(manager *proxy.Manager, adapterId, name string) *VirtualAdapter {
	v := &VirtualAdapter{
		proxy.NewAdapter(manager, adapterId, name),
	}
	v.StartPairing(300 * time.Duration(time.Millisecond))
	return v
}

func (a *VirtualAdapter) StartPairing(timeout time.Duration) {

	discover := func(done <-chan time.Time) {
		for {
			select {
			case <-done:
				fmt.Print("pairing timeout")
				return
			default:
				bulb, err := yeelight.Discover()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				bulb.TurnOff()
				_, err = bulb.SetName("Yeelight" + uuid.NewV4().String())

				params, err := bulb.Discover()
				if err != nil {
					fmt.Printf(err.Error())
					continue
				}
				device := NewYeelightBulb(a, bulb, params)
				a.HandleDeviceAdded(device)
			}
		}
	}
	done := time.After(timeout)
	go discover(done)
}
