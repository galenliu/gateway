package homekit

import (
	"github.com/brutella/hc/characteristic"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/schemas"
	"net"
)

type Handler interface {
	SetValue(name string, value any)
}

type Property struct {
	handler Handler
	*properties.Property
	*characteristic.Characteristic
}

//AlarmProperty                    = "AlarmProperty"
//BarometricPressureProperty       = "BarometricPressureProperty"
//Boolean                  = "Boolean"
//BrightnessProperty               = "BrightnessProperty"
//ColorModeProperty                = "ColorModeProperty"
//ColorProperty                    = "ColorProperty"
//ColorTemperatureProperty         = "ColorTemperatureProperty"
//ConcentrationProperty            = "ConcentrationProperty"
//CurrentProperty                  = "CurrentProperty"
//DensityProperty                  = "DensityProperty"
//FrequencyProperty                = "FrequencyProperty"
//HeatingCoolingProperty           = "HeatingCoolingProperty"
//HumidityProperty                 = "HumidityProperty"
//ImageProperty                    = "ImageProperty"
//InstantaneousPowerFactorProperty = "InstantaneousPowerFactorProperty"
//InstantaneousPowerProperty       = "InstantaneousPowerProperty"
//LeakProperty                     = "LeakProperty"
//LevelProperty                    = "LevelProperty"
//LockedProperty                   = "LockedProperty"
//MotionProperty                   = "MotionProperty"
//OnOffProperty                    = "OnOffProperty"
//OpenProperty                     = "OpenProperty"
//PushedProperty                   = "PushedProperty"
//SmokeProperty                    = "SmokeProperty"
//TargetTemperatureProperty        = "TargetTemperatureProperty"
//TemperatureProperty              = "TemperatureProperty"
//ThermostatModeProperty           = "ThermostatModeProperty"
//VideoProperty                    = "VideoProperty"
//VoltageProperty                  = "VoltageProperty"

func NewProperty(prop *properties.Property) *Property {
	p := &Property{
		Property: prop,
	}

	return p
}

func (p Property) OnValueChanged(v any) {
	p.handler.SetValue(p.GetName(), v)
}

func (p *Property) CreateCharacteristic() {
	switch p.GetAtType() {
	case schemas.OnOffProperty:
		c := characteristic.NewOn()
		c.OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, newValue, oldValue interface{}) {
			p.OnValueChanged(newValue)
		})
		p.Characteristic = c.Characteristic
		break
	case schemas.BrightnessProperty:
		c := characteristic.NewBrightness()
		c.OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, newValue, oldValue interface{}) {
			p.OnValueChanged(newValue)
		})
		p.Characteristic = c.Characteristic
		break

	case schemas.ColorProperty:
		c := characteristic.NewHue()
		c.OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, newValue, oldValue interface{}) {
			p.OnValueChanged(newValue)
		})
		p.Characteristic = c.Characteristic
		break
	}
}
