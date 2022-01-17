package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type BooleanPropertyInstance interface {
	properties.Entity
	Turn(b bool)
}

type BooleanProperty struct {
	BooleanPropertyInstance
}

func NewBoolean(p BooleanPropertyInstance) *BooleanProperty {
	return &BooleanProperty{p}
}

func (b *BooleanProperty) SetValue(a any) {
	v, ok := a.(bool)
	if ok {
		b.BooleanPropertyInstance.Turn(v)
	}
}
