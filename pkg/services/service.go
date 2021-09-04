package services

type Service interface {
	GetID() string
	OnNewThingAdded([]byte)
	OnPropertyChanged([]byte)
	OnAction([]byte)

	SetPropertyValue(v interface{})
}
