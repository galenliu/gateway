package db

import (
	"database/sql"
	"github.com/galenliu/gateway/pkg/logging"
)

func (s *Store) GetSetting(key string) (value string, err error) {
	err = s.db.QueryRow("SELECT value FROM settings where key = @key", sql.Named("key", key), sql.Named("key", key)).Scan(&value)
	return value, err
}

func (s *Store) SetSetting(key, value string) error {

	logging.Info("set setting key:%s value:%s \t\n", key, value)
	_, err := s.GetSetting(key)
	if err == nil {
		_, e := s.db.Exec(`update settings set value=@value where key=@key`, sql.Named("value", value), sql.Named("key", key))
		return e
	}
	stmt, err := s.db.Prepare("INSERT INTO settings(key, value) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, ee := stmt.Exec(key, value)
	if ee != nil {
		return ee
	}
	id, eee := res.LastInsertId()
	if eee != nil {
		return eee
	}
	logging.Debug("insert data,id:%d , value: %s \t\n", id, value)
	return nil
}
