package services

type Service interface {
	OnNewThingAdded([]byte)
	OnPropertyChanged([]byte)
	OnAction([]byte)
}
