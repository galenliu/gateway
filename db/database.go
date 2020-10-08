package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"os"
	"path"
)

const FileName = "db.sqlite3"

var DB *DataBase
var log *zap.Logger

type DataBase struct {
	*sql.DB
	file string
}

func GetDB() *DataBase {
	return DB
}

func InitDB(dir string) error {

	DB = &DataBase{
		file: path.Join(dir, FileName),
	}
	err := DB.open()
	if err != nil {
		return err
	}
	defer DB.close()
	_, err = DB.Exec(sqlCreateTable())
	return err
}

func ResetDB(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		err = os.Remove(dir)
		return err
	}
	return nil
}

func (db *DataBase) open() error {
	var err error
	db.DB, err = sql.Open("sqlite3", db.file)
	return err
}

func (db *DataBase) close() error {
	err := db.Close()
	return err
}

//查询
func (db *DataBase) get(_sql string) (string, error) {
	_ = db.open()
	defer db.close()
	rows, err := db.Query(_sql)
	if err != nil {
		return "", err
	}
	var value string
	err = rows.Scan(value)
	return value, err
}

//插入
func (db *DataBase) insert(prepare string, k, value string) error {
	db.open()
	defer db.close()
	stmt, err := db.Prepare(prepare)
	if err != nil {
		return err
	}
	_, e := stmt.Exec(k, value)
	return e
}

//更新
func (db *DataBase) update(prepare string, arg ...string) error {
	db.open()
	defer db.close()
	stmt, err := db.Prepare(prepare)
	if err != nil {
		return err
	}
	_, e := stmt.Exec(arg)
	return e
}
