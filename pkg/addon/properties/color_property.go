package properties

import (
	"errors"
	"fmt"
	"image/color"
)

type ColorProperty struct {
	*StringProperty
}

func NewColorProperty(value string, args ...string) *ColorProperty {
	desc := ""
	title := "hue"
	if len(args) > 0 {
		desc = args[0]
	}
	if len(args) > 1 {
		title = args[1]
	}
	c := &ColorProperty{}
	c.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:        "color",
		AtType:      TypeColorProperty,
		Title:       title,
		Type:        TypeString,
		Description: desc,
		ReadOnly:    false,
		Value:       value,
	})
	return c
}

func (c *ColorProperty) SetValue(v string) {
	fmt.Print("set value func not implemented")
}

func HTMLToRGB(in string) (color.RGBA, error) {
	if in[0] == '#' {
		in = in[1:]
	}
	if len(in) != 6 {
		return color.RGBA{}, errors.New("Invalid string length")
	}

	var r, g, b byte
	if n, err := fmt.Sscanf(in, "%2x%2x%2x", &r, &g, &b); err != nil || n != 3 {
		return color.RGBA{}, err
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}
