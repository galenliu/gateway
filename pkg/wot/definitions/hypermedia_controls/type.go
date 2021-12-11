package hypermedia_controls

import (
	"github.com/xiam/to"
	"strings"
)

type UnsignedInt uint
type Integer uint
type Number float64
type URI string
type ArrayOrString string
type Double float64

const (
	TypeNumber  = "number"
	TypeString  = "string"
	TypeInteger = "integer"
	TypeNull    = "null"
	TypeObject  = "object"
	TypeArray   = "array"
	TypeBoolean = "boolean"
)

func ToInteger(v interface{}) Integer {
	return Integer(to.Uint64(v))
}

func ToString(v interface{}) string {
	return to.String(v)
}

func ToBool(v interface{}) bool {
	return to.Bool(v)
}

func ToDouble(v interface{}) Double {
	return Double(to.Float64(v))
}

func ToNumber(v interface{}) Number {
	return Number(to.Float64(v))
}

func ToUnsignedInt(v interface{}) UnsignedInt {
	return UnsignedInt(to.Uint64(v))
}

func (i Integer) Compare(value Integer) int {
	if i < value {
		return 1
	}
	if i == value {
		return 0
	}
	return -1
}

func (i Number) Compare(value Number) int {
	if i < value {
		return 1
	}
	if i == value {
		return 0
	}
	return -1
}

func (i UnsignedInt) Compare(value UnsignedInt) int {
	if i < value {
		return 1
	}
	if i == value {
		return 0
	}
	return -1
}

func (u URI) GetId() string {
	l := strings.Split(string(u), "/")
	return l[len(l)-1]
}

func (u URI) GetURI() string {
	return string(u)
}
