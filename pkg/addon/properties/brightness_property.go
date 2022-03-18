package properties

//type BrightnessPropertyDescription struct {
//	Description string  `json:"description,omitempty"`
//	Value       Integer `json:"value,omitempty"`
//}

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty(value Integer, opts ...Option) *BrightnessProperty {

	b := &BrightnessProperty{}
	b.IntegerProperty = NewIntegerProperty(IntegerPropertyDescription{
		Name:    "level",
		AtType:  TypeBrightnessProperty,
		Title:   "brightness",
		Unit:    UnitPercent,
		Minimum: 0,
		Maximum: 100,
		Value:   value,
	}, opts...)
	return b
}
