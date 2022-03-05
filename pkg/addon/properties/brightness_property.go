package properties

//type BrightnessPropertyDescription struct {
//	Description string  `json:"description,omitempty"`
//	Value       Integer `json:"value,omitempty"`
//}

type BrightnessProperty struct {
	*IntegerProperty
}

func NewBrightnessProperty(value Integer, args ...string) *BrightnessProperty {
	desc := ""
	title := "brightness"
	if len(args) > 0 {
		desc = args[0]
	}
	if len(args) > 1 {
		title = args[1]
	}
	b := &BrightnessProperty{}
	b.IntegerProperty = NewIntegerProperty(IntegerPropertyDescription{
		Name:        "level",
		AtType:      TypeBrightnessProperty,
		Title:       title,
		Unit:        UnitPercent,
		Description: desc,
		Minimum:     0,
		Maximum:     100,
		ReadOnly:    false,
		Value:       value,
	})
	return b
}
