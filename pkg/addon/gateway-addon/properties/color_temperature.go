package properties

const TypeColorTemperatureProperty = "ColorTemperatureProperty"

type ColorTemperatureProperty struct {
	*IntegerProperty
}

func NewColorTemperatureProperty() *ColorTemperatureProperty {
	p := NewIntegerProperty(TypeColorTemperatureProperty)
	p.Type = TypeInteger
	p.Name = ColorTemperature
	//p.SetValue(0)
	//p.SetUnit(UnitKelvin)

	return &ColorTemperatureProperty{p}
}
