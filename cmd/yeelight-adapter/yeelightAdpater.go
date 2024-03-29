package yeelight_adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/galenliu/gateway/cmd/yeelight-adapter/lib"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
	"time"
)

type YeeDevice interface {
	HandleRemoved()
}

type YeelightAdapter struct {
	*proxy.Adapter
	devices map[string]YeeDevice
}

func NewVirtualAdapter(adapterId string) *YeelightAdapter {
	v := &YeelightAdapter{
		proxy.NewAdapter(adapterId, "yeelight"),
		make(map[string]YeeDevice),
	}
	return v
}

func (a *YeelightAdapter) HandleDeviceSaved(msg messages.DeviceSavedNotificationJsonData) {
	fmt.Printf("yeelight-adapter handle device saved deviceId: %s \t\n", msg.DeviceId)
}

func (a *YeelightAdapter) StartPairing(timeout <-chan time.Time) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	devices, err := yeelight.Discovery(ctx)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		log.Fatalln(err)
	}
	if len(devices) == 0 {
		fmt.Printf("没有找到Yeelight\t\n")
		return
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
		if found == nil {
			continue
		}
		a.devices[device.ID] = found
		a.HandleDeviceAdded(found)
	}
}

func (a *YeelightAdapter) HandleDeviceRemoved(device proxy.DeviceProxy) {
	dev := a.devices[device.GetId()]
	dev.HandleRemoved()
	delete(a.devices, device.GetId())
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
