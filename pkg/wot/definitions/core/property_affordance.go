package core

type PropertyAffordance interface {
	IsReadOnly() bool
	GetDefaultValue() interface{}
	UnmarshalJSON(data []byte) error
	MarshalJSON() ([]byte, error)
}
