package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
)

const FileName = "database.sqlite3"

var db *gorm.DB

type Config struct {
	file string
}

func InitDB(dir string) error {

	var err error
	c := &Config{
		file: path.Join(dir, FileName),
	}
	db, err = gorm.Open(sqlite.Open(c.file), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func ResetDB(dir string) {
	_, err := os.Stat(dir)
	if err == nil {
		_ = os.Remove(dir)
	}
}

func GetDB() *gorm.DB {
	return db
}
