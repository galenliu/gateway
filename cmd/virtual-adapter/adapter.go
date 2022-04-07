package virtual_adapter

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	p "github.com/galenliu/gateway/pkg/addon/properties"
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
		on := p.NewColorProperty("#ff0000")
		color := p.NewOnOffProperty(true)
		level := p.NewBrightnessProperty(20)
		light.addProperties(on.Property, color.Property, level.Property)
		devs = append(devs, light)
	}

	{
		light := NewVirtualDevice(devices.NewLightBulb("virtual_light1", devices.WithPin("12345678")))
		on := p.NewColorProperty("#ff0000")
		color := p.NewOnOffProperty(true)
		level := p.NewBrightnessProperty(20)
		light.addProperties(on.Property, color.Property, level.Property)
		devs = append(devs, light)
	}

	{
		levelSwitch := NewVirtualDevice(devices.NewMultiLevelSwitch("virtual_level"))
		level := p.NewLevelProperty(10, 0, 100, p.WithUnit(p.UnitPercent))
		onOff := p.NewOnOffProperty(false)
		levelSwitch.addProperties(level.Property, onOff.Property)
		devs = append(devs, levelSwitch)
	}

	{
		levelSwitch := NewVirtualDevice(devices.NewMultiLevelSwitch("virtual_level1", devices.WithCredentialsRequired()))
		level := p.NewLevelProperty(10, 0, 100)
		onOff := p.NewOnOffProperty(false)
		levelSwitch.addProperties(level.Property, onOff.Property)
		devs = append(devs, levelSwitch)
	}

	{
		smartPlug := NewVirtualDevice(devices.NewSmartPlug("smart_plug"))
		//开关
		onOff := p.NewOnOffProperty(false)
		//级别
		level := p.NewLevelProperty(10, 0, 100)
		//功率
		power := p.NewInstantaneousPowerProperty(0, p.WithReadOnly())
		//频率
		factor := p.NewInstantaneousPowerFactorProperty(0, p.WithReadOnly())
		//电压
		voltage := p.NewVoltageProperty(0, p.WithReadOnly())
		//电流
		current := p.NewCurrentProperty(0, p.WithReadOnly())
		//频率
		frequency := p.NewFrequencyProperty(0, p.WithReadOnly())
		smartPlug.addProperties(onOff.Property, level.Property, power.Property, factor.Property,
			voltage.Property, current.Property, frequency.Property)
		devs = append(devs, smartPlug)
	}

	{
		//运动传感器
		motionSensor := NewVirtualDevice(devices.NewMotionSensor("virtual_motion_sensor"))
		motion := p.NewLevelProperty(0, 0, 100, p.WithReadOnly())
		motionSensor.addProperties(motion.Property)
		devs = append(devs, motionSensor)
	}

	{
		//多级传感器
		motionLevelSensor := NewVirtualDevice(devices.NewMultiLevelSensor("virtual_motion_level_sensor"))
		motion := p.NewLevelProperty(0, 0, 100, p.WithReadOnly())
		motionLevelSensor.addProperties(motion.Property)
		devs = append(devs, motionLevelSensor)
	}

	{
		//空调
		thermostat := NewVirtualDevice(devices.NewThermostat("virtual_thermostat"))
		//当前温度
		temperature := p.NewTemperatureProperty(20, p.WithReadOnly())
		//制热目标温度
		coolingTargetTemperature := p.NewTargetTemperatureProperty(25, p.WithTitle("Cooling Target"), p.WithName("coolingTargetTemperature"))
		//制冷目标温度
		heatingTargetTemperature := p.NewTargetTemperatureProperty(19, p.WithTitle("Heating Target"), p.WithName("heatingTargetTemperature"))

		//制热制冷属性
		heatingCooling := p.NewHeatingCoolingProperty(p.ThermostatCool, p.WithTitle("Heating/Cooling"))

		//制热制冷模式
		thermostatMode := p.NewThermostatModeProperty(p.ThermostatCool)

		thermostat.addProperties(temperature.Property, heatingTargetTemperature.Property, coolingTargetTemperature.Property, heatingCooling.Property, thermostatMode.Property)
		devs = append(devs, thermostat)
	}
	a.HandleDeviceAdded(devs...)

}

func (a *Adapter) HandleDeviceSaved(msg messages.DeviceSavedNotificationJsonData) {
	fmt.Printf("virtual-adapter handle device saved deviceId: %s \t\n", msg.DeviceId)
}
