package properties

type HeatingCoolingEnum string

var HeatingCoolingEnumOff HeatingCoolingEnum = "off"
var HeatingCoolingEnumHeating HeatingCoolingEnum = "heating"
var HeatingCoolingEnumCooling HeatingCoolingEnum = "cooling"
var HeatingCoolingEnumAuto HeatingCoolingEnum = "auto"

type HeatingCoolingProperty struct {
	*StringProperty
}

func NewHeatingCoolingProperty(value HeatingCoolingEnum, enums []HeatingCoolingEnum, opts ...Option) *HeatingCoolingProperty {
	b := &HeatingCoolingProperty{}
	opts = append(opts, WithUnit(UnitAmpere))
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name: "heatingCooling",
		Enum: func() []string {
			enum := make([]string, 0)
			for _, e := range enums {
				enum = append(enum, string(e))
			}
			return enum
		}(),
		AtType: TypeHeatingCoolingProperty,
		Value:  string(value),
	}, opts...)
	return b
}
