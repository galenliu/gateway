package db

func (s *Storage) LoadServiceSetting(key string) (value string, err error) {
	return s.GetSetting("service." + key)
}

func (s *Storage) StoreServiceSetting(key, value string) error {
	return s.SetSetting("service."+key, value)
}

func (s *Storage) LoadServiceConfig(key string) (value string, err error) {
	return s.GetSetting("service.config" + key)
}

func (s *Storage) StoreServiceConfig(key, value string) error {
	return s.SetSetting("service.config"+key, value)
}
