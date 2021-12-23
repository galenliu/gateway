package core

type PropertyAffordance interface {
	IsReadOnly() bool
	GetDefaultValue() any
}
