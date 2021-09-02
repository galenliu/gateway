package db

func (s *Storage) LoadAddonSetting(key string) (value string, err error) {
	return s.GetSetting("addons." + key)
}

func (s *Storage) StoreAddonSetting(key, value string) error {
	return s.SetSetting("addons."+key, value)
}

func (s *Storage) LoadAddonConfig(key string) (value string, err error) {
	return s.GetSetting("addons.config" + key)
}

func (s *Storage) StoreAddonsConfig(key, value string) error {
	return s.SetSetting("addons.config"+key, value)
}
