package virtual_adapter

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
		light := NewVirtualDevice(devices.NewLightBulb("virtual_light"))
		on := properties.NewColorProperty("#ff0000")
		color := properties.NewOnOffProperty(true)
		level := properties.NewBrightnessProperty(20)
		light.addProperties(on.Property, color.Property, level.Property)
		devs = append(devs, light)
	}

	{
		light := NewVirtualDevice(devices.NewLightBulb("virtual_light1", devices.WithPin("12345678")))
		on := properties.NewColorProperty("#ff0000")
		color := properties.NewOnOffProperty(true)
		level := properties.NewBrightnessProperty(20)
		light.addProperties(on.Property, color.Property, level.Property)
		devs = append(devs, light)
	}

	{
		levelSwitch := NewVirtualDevice(devices.NewMultiLevelSwitch("virtual_level"))
		level := properties.NewLevelProperty(10, 0, 100, properties.WithUnit(properties.UnitPercent))
		onOff := properties.NewOnOffProperty(false)
		levelSwitch.addProperties(level.Property, onOff.Property)
		devs = append(devs, levelSwitch)
	}

	{
		levelSwitch := NewVirtualDevice(devices.NewMultiLevelSwitch("virtual_level1", devices.WithCredentialsRequired()))
		level := properties.NewLevelProperty(10, 0, 100)
		onOff := properties.NewOnOffProperty(false)
		levelSwitch.addProperties(level.Property, onOff.Property)
		devs = append(devs, levelSwitch)
	}

	{
		smartPlug := NewVirtualDevice(devices.NewSmartPlug("smart_plug"))
		//开关
		onOff := properties.NewOnOffProperty(false)
		//级别
		level := properties.NewLevelProperty(10, 0, 100)
		//功率
		power := properties.NewInstantaneousPowerProperty(0, properties.WithReadOnly())
		//频率
		factor := properties.NewInstantaneousPowerFactorProperty(0, properties.WithReadOnly())
		//电压
		voltage := properties.NewVoltageProperty(0, properties.WithReadOnly())
		//电流
		current := properties.NewCurrentProperty(0, properties.WithReadOnly())
		//频率
		frequency := properties.NewFrequencyProperty(0, properties.WithReadOnly())
		smartPlug.addProperties(onOff.Property, level.Property, power.Property, factor.Property,
			voltage.Property, current.Property, frequency.Property)
		devs = append(devs, smartPlug)
	}

	{
		//运动传感器
		motionSensor := NewVirtualDevice(devices.NewMotionSensor("virtual_motion_sensor"))
		motion := properties.NewLevelProperty(0, 0, 100, properties.WithReadOnly())
		motionSensor.addProperties(motion.Property)
		devs = append(devs, motionSensor)

	}

	{
		//多级传感器
		motionLevelSensor := NewVirtualDevice(devices.NewMultiLevelSensor("virtual_motion_level_sensor"))
		motion := properties.NewLevelProperty(0, 0, 100, properties.WithReadOnly())
		motionLevelSensor.addProperties(motion.Property)
		devs = append(devs, motionLevelSensor)
	}

	{
		//空调
		thermostat := NewVirtualDevice(devices.NewThermostat("virtual_thermostat"))
		temperature := properties.NewTemperatureProperty(0, properties.WithReadOnly())
		targetTemperature := properties.NewTargetTemperatureProperty(25)
		heatingCooling := properties.NewHeatingCoolingProperty(properties.HeatingCoolingEnumOff)
		thermostatMode := properties.NewThermostatModeProperty(properties.ThermostatModeEnumAuto)
		thermostat.addProperties(temperature.Property, targetTemperature.Property, heatingCooling.Property, thermostatMode.Property)
		devs = append(devs, thermostat)
	}
	a.HandleDeviceAdded(devs...)

}

func (a *Adapter) HandleDeviceSaved(msg messages.DeviceSavedNotificationJsonData) {
	fmt.Printf("virtual-adapter handle device saved deviceId: %s \t\n", msg.DeviceId)
}
