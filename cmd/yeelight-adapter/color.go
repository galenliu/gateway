package yeelight_adapter

import (
	"context"
	"github.com/galenliu/gateway/cmd/yeelight-adapter/lib"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type Color struct {
	device *YeelightDevice
	*properties.ColorProperty
}

func NewColor(bulb *YeelightDevice, value string) *Color {
	return &Color{
		bulb,
		properties.NewColorProperty(value),
	}
}

func (on *Color) SetValue(v string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	rgb := strings.TrimPrefix(v, "#")
	color64, err := strconv.ParseInt(rgb, 16, 32)
	if err != nil {
		return err
	}
	err = on.device.Client.SetRGB(ctx, int(color64), yeelight.EffectSudden, duration)
	if err != nil {
		return err
	}
	on.SetCachedValue(v)
	on.NotifyChanged()
	return nil
}

func RGBToYeelight(color color.RGBA) int {
	r := int(color.R)
	g := int(color.G)
	b := int(color.B)

	return r*65536 + g*256 + b
}
