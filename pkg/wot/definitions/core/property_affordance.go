package core

type PropertyAffordance interface {
	IsReadOnly() bool
	GetDefaultValue() any
	UnmarshalJSON(data []byte) error
	MarshalJSON() ([]byte, error)
}
