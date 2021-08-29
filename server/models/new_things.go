package models

type NewThingsManager interface {
	GetDevicesBytes() map[string][]byte
}

type NewThingsModel struct {
	Manager NewThingsManager
}

func NewNewThings() *NewThingsModel {
	n := &NewThingsModel{}
	return n
}
