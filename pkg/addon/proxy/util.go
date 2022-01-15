package proxy

import (
	"errors"
	"fmt"
	"image/color"
)

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
	return color.RGBA{R: uint8(r) / 255, G: uint8(g) / 255, B: uint8(b) / 255}, nil
}