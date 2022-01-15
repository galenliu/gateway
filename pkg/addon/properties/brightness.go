package properties

const TypeBrightnessProperty = "BrightnessProperty"

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty(description PropertyDescription) *BrightnessProperty {
	p := &IntegerProperty{}
	p.Property = NewProperty(description)
	if description.Name == nil {
		p.Name = "bright"
	}
	if description.Maximum == nil {
		p.SetMinValue(100)
	}
	if description.Minimum == nil {
		p.SetMinValue(0)
	}
	if description.Unit == nil {
		p.Unit = UnitPercentage
	}
	return &BrightnessProperty{p}
}

func (b BrightnessProperty) SetBrightness(v int) {

}
