package virtual

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type Adapter struct {
	*proxy.Adapter
}

func NewVirtualAdapter(adapterId string) *Adapter {
	v := &Adapter{
		proxy.NewAdapter(adapterId, "Virtual"),
	}
	return v
}

func (a *Adapter) StartPairing(t <-chan time.Time) {
	devs := make([]proxy.DeviceProxy, 0)

	{
		light := NewVirtualDevice(devices.NewLightBulb("light1"))
		on := properties.NewColorProperty("#ff0000")
		color := properties.NewOnOffProperty(true)
		level := properties.NewBrightnessProperty(20)
		light.addProperties(on.Property, color.Property, level.Property)
		devs = append(devs, light)
	}

	{
		levelSwitch := NewVirtualDevice(devices.NewMultiLevelSwitch("cover1"))
		level := properties.NewLevelProperty(10, 0, 100, properties.WithUnit(properties.UnitPercent))
		onOff := properties.NewOnOffProperty(false)
		levelSwitch.addProperties(level.Property, onOff.Property)
		devs = append(devs, levelSwitch)
	}

	{
		light := NewVirtualDevice(devices.NewLightBulb("light2", devices.WithPin("12345678")))
		on := properties.NewColorProperty("#ff0000")
		color := properties.NewOnOffProperty(true)
		level := properties.NewBrightnessProperty(20)
		light.addProperties(on.Property, color.Property, level.Property)
		devs = append(devs, light)
	}

	{
		levelSwitch := NewVirtualDevice(devices.NewMultiLevelSwitch("cover2", devices.WithCredentialsRequired()))
		level := properties.NewLevelProperty(10, 0, 100)
		onOff := properties.NewOnOffProperty(false)
		levelSwitch.addProperties(level.Property, onOff.Property)
		devs = append(devs, levelSwitch)
	}

	a.AddDevices(devs...)
}

func (a *Adapter) HandleDeviceSaved(msg messages.DeviceSavedNotificationJsonData) {
	fmt.Printf("virtual-adapter handle device saved deviceId: %s \t\n", msg.DeviceId)
}
