package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
)

type BooleanValuer interface {
	devices.PropertyEntity
	SetValue(bool2 bool)
}

type Boolean struct {
	BooleanValuer
}

func NewBoolean(b BooleanValuer) *Boolean {
	return &Boolean{b}
}

func (b Boolean) SetValue(a any) {
	v, ok := a.(bool)
	if ok {
		b.BooleanValuer.SetValue(v)
	}

}
