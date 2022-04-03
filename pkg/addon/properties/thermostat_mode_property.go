package properties

type ThermostatModeEnum string

var ThermostatModeEnumOff = "off"
var ThermostatModeEnumHeat = "heat"
var ThermostatModeEnumCool = "cool"
var ThermostatModeEnumAuto = "auto"

type ThermostatModeProperty struct {
	*StringProperty
}

func NewThermostatModeProperty(value string, opts ...Option) *ThermostatModeProperty {
	b := &ThermostatModeProperty{}
	opts = append(opts, WithTitle("ThermostatMode"), WithUnit(UnitVolt))
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:     "thermostatMode",
		ReadOnly: true,
		Enum:     []string{ThermostatModeEnumOff, ThermostatModeEnumHeat, ThermostatModeEnumCool, ThermostatModeEnumAuto},
		AtType:   TypeThermostatModeProperty,
		Value:    value,
	}, opts...)
	return b
}
