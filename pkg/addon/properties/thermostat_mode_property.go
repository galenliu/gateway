package properties

type ThermostatModeEnum string

var ThermostatModeEnumOff ThermostatModeEnum = "off"
var ThermostatModeEnumHeat ThermostatModeEnum = "heat"
var ThermostatModeEnumCool ThermostatModeEnum = "cool"
var ThermostatModeEnumAuto ThermostatModeEnum = "auto"

type ThermostatModeProperty struct {
	*StringProperty
}

func NewThermostatModeProperty(value ThermostatModeEnum, opts ...Option) *ThermostatModeProperty {
	b := &ThermostatModeProperty{}
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:   "thermostatMode",
		Enum:   []string{string(ThermostatModeEnumOff), string(ThermostatModeEnumHeat), string(ThermostatModeEnumCool), string(ThermostatModeEnumAuto)},
		Unit:   UnitDegreeCelsius,
		Title:  "Mode",
		AtType: TypeThermostatModeProperty,
		Value:  string(value),
	}, opts...)
	return b
}
