package db

const SettingsTable = "settings"

func (s *Storage) GetSetting(key string) (value string, err error) {
	return s.queryValue(key, SettingsTable)
}

func (s *Storage) SetSetting(key, value string) error {
	return s.setValue(key, value, SettingsTable)
}

func (s *Storage) RemoveSetting(key string) error {
	return s.deleteValue(key, SettingsTable)
}
