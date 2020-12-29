package database

import (
	"gateway/pkg/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
)

const FileName = "database.sqlite3"

var db *gorm.DB

var config *Config

type Config struct {
	file string
}

func InitDB(dir string) error {

	var err error
	config := &Config{
		file: path.Join(dir, FileName),
	}
	db, err = gorm.Open(sqlite.Open(config.file), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Debug("database init succeed")
	return nil
}

func ResetDB(dir string) {
	_, err := os.Stat(dir)
	if err == nil {
		_ = os.Remove(dir)
	}
}

func GetDB() (*gorm.DB, error) {
	var err error
	if db != nil {
		return db, nil
	}
	db, err = gorm.Open(sqlite.Open(config.file), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}
