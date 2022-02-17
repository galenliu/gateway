package properties

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty(description PropertyDescription) *BrightnessProperty {
	p := &BrightnessProperty{}
	description.AtType = "BrightnessProperty"
	if description.Maximum == nil {
		description.Maximum = 100
	}
	if description.Minimum == nil {
		description.Minimum = 0
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
