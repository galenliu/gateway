package properties

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty(description PropertyDescription) *BrightnessProperty {
	p := &BrightnessProperty{}
	description.AtType = "BrightnessProperty"
	if description.Maximum == nil {
		var v float64 = 100
		description.Maximum = &v
	}
	if description.Minimum == nil {
		var v float64 = 0
		description.Minimum = &v
	}
	if description.Unit == "" {
		description.Unit = "percentage"
	}
	if description.Name == "" {
		description.Name = "bright"
	}
	p.IntegerProperty = NewIntegerProperty(description)
	return p
}
