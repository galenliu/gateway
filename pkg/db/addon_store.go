package db

func (s *Storage) GetAddonsSetting(key string) (value string, err error) {
	return s.GetSetting("addons." + key)
}

func (s *Storage) SetAddonsSetting(key, value string) error {
	return s.SetSetting("addons."+key, value)
}

func (s *Storage) GetAddonsConfig(key string) (value string, err error) {
	return s.GetSetting("addons.config" + key)
}

func (s *Storage) SetAddonsConfig(key, value string) error {
	return s.SetSetting("addons.config"+key, value)
}
