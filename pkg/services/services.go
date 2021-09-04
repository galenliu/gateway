package services

type ThingsManager interface {
	SetPropertyValue(thingId, propertyName string, value interface{}) (interface{}, error)
	GetPropertyValue(thingId, propertyName string) (interface{}, error)
	GetPropertiesValue(thingId string) (map[string]interface{}, error)
}

type Services struct {
	services map[string]Service
	manager  ThingsManager
}

func NewServices(m ThingsManager) *Services {
	s := &Services{}
	s.manager = m
	return s
}

func (s *Services) AddService(ser Service) {
	s.services[ser.GetID()] = ser
}

func (s *Services) RemoveService(ser Service) {
	delete(s.services, ser.GetID())
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
