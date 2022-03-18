package yeelight

import (
	"context"
	"errors"
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/lib"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	devices, err := yeelight.Discovery(ctx)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		log.Fatalln(err)
	}

	for _, device := range devices {
		fmt.Println(`------`)
		fmt.Printf("Device '%s' [ID:%s Version:%s]\n", device.Name, device.ID, device.FirmwareVersion)
		fmt.Printf("Address: %s\n", device.Location)
		fmt.Printf("Power: %s\n", power(device.Power))
		fmt.Printf("Support: %s\n", device.Support)
		fmt.Printf("Model: %s\n", device.Model)

		// create new client for work with device
		client := yeelight.New(device.Location)
		if device.Name == "" || device.Name != "yeelight"+device.Location {
			_ = client.SetName(context.Background(), "yeelight"+device.Location)
		}
		found := NewYeelightBulb(&client, device.ID, device.Name, device.Location)
		a.HandleDeviceAdded(found)

	}
	//fmt.Printf("start pairing...\n")
	//if timeout == nil {
	//	timeout = time.After(5 * time.Second)
	//}
	//
	//discover := func() *yeelight.Yeelight {
	//	bulb, err := yeelight.Discover()
	//	if err != nil {
	//		fmt.Printf("adapter:%s  err:%s \t\n", a.GetId(), err.Error())
	//		return nil
	//	}
	//	return bulb
	//}
	//for {
	//	foundDevice := discover()
	//	if foundDevice != nil {
	//		yl := NewYeelightBulb(foundDevice)
	//		a.HandleDeviceAdded(yl)
	//	}
	//	select {
	//	case <-timeout:
	//		return
	//	}
	//}

}

func (a *YeelightAdapter) HandleDeviceRemoved(device proxy.DeviceProxy) {

}

func power(on bool) string {
	if on {
		return "ON"
	}
	return "OFF"
}

func isPowerOn(power string) bool {
	if power == "on" {
		return true
	}
	return false
}
