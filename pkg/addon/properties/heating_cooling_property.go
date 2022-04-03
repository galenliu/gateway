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
	opts = append(opts, WithTitle("HeatingCooling"), WithUnit(UnitAmpere))
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:   "heating_cooling",
		Enum:   []string{string(HeatingCoolingEnumOff), string(HeatingCoolingEnumHeating), string(HeatingCoolingEnumCooling)},
		AtType: TypeHeatingCoolingProperty,
		Value:  string(value),
	}, opts...)
	return b
}
