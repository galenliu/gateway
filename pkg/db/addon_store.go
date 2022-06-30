package db

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

const AddonTable = "settings"

func (s *Storage) LoadAddonSetting(key string) (value string, err error) {
	return s.GetSetting("addons." + key)
}

func (s *Storage) UpdateAddonSetting(id, value string) (err error) {
	key := "addons." + id
	return s.updateValue(key, value, AddonTable)
}

func (s *Storage) UpdateAddonConfig(id, value string) (err error) {
	key := "addons.config." + id
	return s.updateValue(key, value, AddonTable)
}

func (s *Storage) StoreAddonSetting(key, value string) error {
	return s.setValue("addons."+key, value, AddonTable)
}

func (s *Storage) LoadAddonConfig(key string) (value string, err error) {
	return s.GetSetting("addons.config." + key)
}

func (s *Storage) StoreAddonsConfig(key string, v any) error {
	value, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.SetSetting("addons.config."+key, gjson.GetBytes(value, "@pretty").String())
}

func (s *Storage) RemoveAddonSettingAndConfig(key string) error {
	err := s.deleteValue("addons.config."+key, AddonTable)
	err1 := s.deleteValue("addons."+key, AddonTable)
	if err != nil {
		return err
	}
	if err1 != nil {
		return err1
	}
	return nil
}
