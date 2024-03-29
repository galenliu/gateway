package hypermedia_controls

import (
	"github.com/xiam/to"
	"strings"
)

type UnsignedInt uint
type Integer uint
type Number float64
type URI string
type ArrayOrString []string
type Double float64

const DefaultInteger Integer = 0
const DefaultNumber Number = 0
const DefaultDouble Double = 0

type DataSchemaType = string

const (
	TypeNumber  DataSchemaType = "number"
	TypeString  DataSchemaType = "string"
	TypeInteger DataSchemaType = "integer"
	TypeNull    DataSchemaType = "null"
	TypeObject  DataSchemaType = "object"
	TypeArray   DataSchemaType = "array"
	TypeBoolean DataSchemaType = "boolean"
)

func (u URI) ToString() string {
	return string(u)
}

func NewURI(s string) URI {
	return URI(s)
}

func NewArrayOrString(args ...string) ArrayOrString {
	arr := make([]string, 0)
	for _, s := range args {
		arr = append(arr, s)
	}
	return arr
}

func ToInteger(v any) Integer {
	return Integer(to.Uint64(v))
}

func ToString(v any) string {
	return to.String(v)
}

func ToBool(v any) bool {
	return to.Bool(v)
}

func ToDouble(v any) Double {
	return Double(to.Float64(v))
}

func ToNumber(v any) Number {
	return Number(to.Float64(v))
}

func ToUnsignedInt(v any) UnsignedInt {
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
