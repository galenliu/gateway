package db

import (
	"fmt"
)

var SettingsTableName = "settings"
var RowKey = "key"
var RowValue = "value"

type Settings interface {
	Get(key string) string
	Set(key string, value string)
}

func (db *DataBase) SettingGet(key string) (string, error) {
	value, err := db.get(sqlSelectValueByKey(key))
	if err != nil {
		return "", err
	}
	return value, nil
}

func (db *DataBase) SettingSet(key string, value string) error {
	err := db.insert(sqlInsertValueForKey(), key, value)
	return err
}

func (db *DataBase) SettingUpdate(key string, value string) error {
	err := db.update(sqlUpdateForKey(key), value)
	return err
}

func sqlCreateTable() string {
	q := fmt.Sprintf(`
	create table if not exists %v(%v TEXT PRIMARY KEY, %v TEXT);
	delete from settings;
 `, SettingsTableName, RowKey, RowValue)
	return q
}

func sqlSelectValueByKey(key string) string {
	return "SELECT " + RowValue + " FROM " + SettingsTableName + " WHERE " + RowValue + "=" + key
}

func sqlInsertValueForKey() string {
	return "insert into " + SettingsTableName + "(" + RowKey + "," + RowValue + ") values(?,?)"
}

func sqlUpdateForKey(key string) string {
	return "update userinfo set " + RowValue + "=? where " + RowKey + "=" + key
}
