package hypermedia_controls

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
