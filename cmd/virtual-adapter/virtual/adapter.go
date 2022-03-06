package virtual

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/proxy"
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
		light := NewDevice(devices.NewLightBulb("light1"))
		on := properties.NewColorProperty("#ffffff")
		color := properties.NewOnOffProperty(true)
		level := properties.NewBrightnessProperty(20)
		light.addProperties(on, color, level)
		devs = append(devs, light)
	}

	{
		levelSwitch := NewDevice(devices.NewMultiLevelSwitch("cover1"))
		level := properties.NewLevelProperty(10, 0, 100, properties.WithUnit(properties.UnitPercent))
		onOff := properties.NewOnOffProperty(false)
		levelSwitch.addProperties(level, onOff)
		devs = append(devs, levelSwitch)
	}

	{
		light := NewDevice(devices.NewLightBulb("light2", devices.WithPin("12345678")))
		on := properties.NewColorProperty("#ffffff")
		color := properties.NewOnOffProperty(true)
		level := properties.NewBrightnessProperty(20)
		light.addProperties(on, color, level)
		devs = append(devs, light)
	}

	{
		levelSwitch := NewDevice(devices.NewMultiLevelSwitch("cover2", devices.WithCredentialsRequired()))
		level := properties.NewLevelProperty(10, 0, 100)
		onOff := properties.NewOnOffProperty(false)
		levelSwitch.addProperties(level, onOff)
		devs = append(devs, levelSwitch)
	}

	a.AddDevices(devs...)
}
