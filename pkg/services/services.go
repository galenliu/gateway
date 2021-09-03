package services

type Services struct {
	services map[string]Service
}

func NewServices() *Services {
	return nil
}

func (s *Services) NewThingAdded(data []byte) {
	for _, ser := range s.GetServices() {
		ser.OnNewThingAdded(data)
	}
}

func (s *Services) PropertyChanged(data []byte) {
	for _, ser := range s.GetServices() {
		ser.OnPropertyChanged(data)
	}
}

func (s Services) NotifyAction(data []byte) {
	for _, ser := range s.GetServices() {
		ser.OnAction(data)
	}
}

func (s *Services) GetServices() []Service {
	return nil
}
