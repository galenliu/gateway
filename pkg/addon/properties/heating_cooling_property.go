package properties

type HeatingCoolingEnum string

var HeatingCoolingEnumOff HeatingCoolingEnum = "off"
var HeatingCoolingEnumHeating HeatingCoolingEnum = "heating"
var HeatingCoolingEnumCooling HeatingCoolingEnum = "cooling"

type HeatingCoolingProperty struct {
	*StringProperty
}

func NewHeatingCoolingProperty(value HeatingCoolingEnum, opts ...Option) *HeatingCoolingProperty {
	b := &HeatingCoolingProperty{}
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:   "heatingCooling",
		Enum:   []string{"off", "heating", "cooling"},
		Title:  "HeatingCooling",
		Unit:   UnitDegreeCelsius,
		AtType: TypeHeatingCoolingProperty,
		Value:  string(value),
	}, opts...)
	return b
}
