package hypermedia_controls

import "strings"

type UnsignedInt uint
type Integer uint
type Number float64
type URI string

const (
	TypeNumber  = "number"
	TypeString  = "string"
	TypeInteger = "integer"
	TypeNull    = "null"
	TypeObject  = "object"
	TypeArray   = "array"
	TypeBoolean = "boolean"
)

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

func (u URI) GetID() string {
	l := strings.Split(string(u), "/")
	return l[len(l)-1]
}

func (u URI) GetURI() string {
	return string(u)
}
