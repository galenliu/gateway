package properties

const TypeBrightnessProperty = "BrightnessProperty"

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty(description PropertyDescription) *BrightnessProperty {
	p := &BrightnessProperty{}
	atType := "BrightnessProperty"
	description.AtType = &atType
	if description.Maximum == nil {
		var v float64 = 100
		description.Maximum = &v
	}
	if description.Minimum == nil {
		var v float64 = 0
		description.Minimum = &v
	}
	if description.Unit == nil {
		var u = "percentage"
		description.Unit = &u
	}
	if description.Name == nil {
		var name = "bright"
		description.Name = &name
	}
	p.IntegerProperty = NewIntegerProperty(description)
	return p
}

func (b BrightnessProperty) SetBrightness(v int) {

}
