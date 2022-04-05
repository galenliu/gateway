package properties

type ThermostatModeEnum string

var ThermostatModeEnumOff ThermostatModeEnum = "off"
var ThermostatModeEnumHeat ThermostatModeEnum = "heat"
var ThermostatModeEnumCool ThermostatModeEnum = "cool"
var ThermostatModeEnumAuto ThermostatModeEnum = "auto"

type ThermostatModeProperty struct {
	*StringProperty
}

func NewThermostatModeProperty(value ThermostatModeEnum, enum []ThermostatModeEnum, opts ...Option) *ThermostatModeProperty {
	b := &ThermostatModeProperty{}
	opts = append(opts, WithTitle("Mode"), WithUnit(UnitVolt))
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name: "thermostatMode",
		Enum: func() []string {
			em := make([]string, 0)
			for _, e := range enum {
				em = append(em, string(e))
			}
			return em
		}(),
		AtType: TypeThermostatModeProperty,
		Value:  string(value),
	}, opts...)
	return b
}
