package properties

type ThermostatModeEnum string

var ThermostatOff ThermostatModeEnum = "off"
var ThermostatHeat ThermostatModeEnum = "heat"
var ThermostatCool ThermostatModeEnum = "cool"
var ThermostatAuto ThermostatModeEnum = "auto"
var ThermostatDry ThermostatModeEnum = "dry"
var ThermostatWind ThermostatModeEnum = "wind"

type ThermostatModeProperty struct {
	*StringProperty
}

func NewThermostatModeProperty(value ThermostatModeEnum, opts ...Option) *ThermostatModeProperty {
	b := &ThermostatModeProperty{}
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:   "thermostatMode",
		Enum:   []string{string(ThermostatOff), string(ThermostatHeat), string(ThermostatCool), string(ThermostatAuto), string(ThermostatDry), string(ThermostatWind)},
		Unit:   UnitDegreeCelsius,
		Title:  "Mode",
		AtType: TypeThermostatModeProperty,
		Value:  string(value),
	}, opts...)
	return b
}
