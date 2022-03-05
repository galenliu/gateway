package properties

type LevelProperty struct {
	*NumberProperty
}

func NewLevelProperty(value Number, min, max Number, unit Unit, args ...string) *LevelProperty {
	desc := ""
	title := "brightness"
	if len(args) > 0 {
		desc = args[0]
	}
	if len(args) > 1 {
		title = args[1]
	}
	l := &LevelProperty{}
	l.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:        "level",
		AtType:      TypeLevelProperty,
		Title:       title,
		Unit:        unit,
		Description: desc,
		Minimum:     min,
		Maximum:     max,
		ReadOnly:    false,
		Value:       value,
	})
	return l
}
