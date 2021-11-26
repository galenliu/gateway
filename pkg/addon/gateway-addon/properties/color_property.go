package properties

const TypeColorProperty = "ColorProperty"

type ColorProperty struct {
	*StringProperty
}

func NewColorProperty() *ColorProperty {
	p := NewStringProperty(TypeColorProperty)
	p.Type = TypeString
	p.Name = ColorModel
	p.SetValue("#121212")
	p.Unit = UnitPercentage

	return &ColorProperty{p}
}
