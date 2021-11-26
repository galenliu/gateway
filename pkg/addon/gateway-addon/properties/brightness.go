package properties

const TypeBrightnessProperty = "BrightnessProperty"

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty() *BrightnessProperty {
	p := &IntegerProperty{}
	p.Property.Type = TypeInteger
	p.Name = "bright"
	p.SetMinValue(0)
	p.SetMaxValue(100)
	//p.SetValue(0)
	//p.SetUnit(UnitPercentage)

	return &BrightnessProperty{p}
}
