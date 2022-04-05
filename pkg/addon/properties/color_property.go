package properties

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type ColorProperty struct {
	*StringProperty
}

func NewColorProperty(value string, opts ...Option) *ColorProperty {
	c := &ColorProperty{}
	c.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:     "color",
		AtType:   TypeColorProperty,
		ReadOnly: false,
		Value:    value,
	}, opts...)
	return c
}

func (c *ColorProperty) SetValue(v string) {
	fmt.Print("set value func not implemented")
}

type RGB struct {
	R, G, B float64
}

// Takes a string like '#123456' or 'ABCDEF' and returns an RGB
func HTMLToRGB(in string) (RGB, error) {

	in = strings.ToUpper(in)
	if in[0] == '#' {
		in = in[1:]
	}

	if len(in) != 6 {
		return RGB{}, errors.New("Invalid string length")
	}

	var r, g, b byte
	if n, err := fmt.Sscanf(in, "%2x%2x%2x", &r, &g, &b); err != nil || n != 3 {
		return RGB{}, err
	}

	return RGB{float64(r), float64(g), float64(b)}, nil
}

func (c RGB) ToHSL() HSL {
	var h, s, l float64

	r := c.R
	g := c.G
	b := c.B

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)

	// Luminosity is the average of the max and min rgb color intensities.
	l = (max + min) / 2

	// saturation
	delta := max - min
	if delta == 0 {
		// it's gray
		return HSL{0, 0, l}
	}

	// it's not gray
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}

	// hue
	r2 := (((max - r) / 6) + (delta / 2)) / delta
	g2 := (((max - g) / 6) + (delta / 2)) / delta
	b2 := (((max - b) / 6) + (delta / 2)) / delta
	switch {
	case r == max:
		h = b2 - g2
	case g == max:
		h = (1.0 / 3.0) + r2 - b2
	case b == max:
		h = (2.0 / 3.0) + g2 - r2
	}

	// fix wraparounds
	switch {
	case h < 0:
		h += 1
	case h > 1:
		h -= 1
	}

	return HSL{h, s, l}
}

// A nudge to make truncation round to nearest number instead of flooring
const delta = 1 / 512.0

func (c RGB) ToHTML() string {
	return fmt.Sprintf("%02x%02x%02x", byte((c.R+delta)*255), byte((c.G+delta)*255), byte((c.B+delta)*255))
}

type HSL struct {
	H, S, L float64
}

func hueToRGB(v1, v2, h float64) float64 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}

func (c HSL) ToRGB() RGB {
	h := c.H
	s := c.S
	l := c.L

	if s == 0 {
		// it's gray
		return RGB{l, l, l}
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return RGB{r, g, b}
}

func (c HSL) ToHTML() string {
	return c.ToRGB().ToHTML()
}
