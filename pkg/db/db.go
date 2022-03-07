package db

import (
	"database/sql"
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

const dbDefaultFileName = "database.sqlite3"

type Config struct {
	Reset    bool
	FileName string
}

type Storage struct {
	file   string
	db     *sql.DB
	logger logging.Logger
}

func NewStorage(filePath string, log logging.Logger, conf ...Config) (*Storage, error) {
	s := &Storage{}
	s.logger = log
	var config = Config{
		Reset: false,
	}
	if len(conf) > 0 {
		config = conf[0]
	}
	if config.FileName == "" {
		config.FileName = dbDefaultFileName
	}
	s.file = path.Join(filePath, config.FileName)
	if config.Reset {
		s.reset()
	}
	d, e := sql.Open("sqlite3", s.file)
	if e != nil {
		return nil, e
	}
	s.db = d
	ee := s.createTable()
	if ee != nil {
		return nil, ee
	}

	return s, nil
}

func (s *Storage) reset() {
	_, err := os.Stat(s.file)
	if err == nil {
		_ = os.Remove(s.file)
	}
}

func (s *Storage) createTable() error {

	thingsTable := `
    CREATE TABLE IF NOT EXISTS things(
        id TEXT PRIMARY KEY NOT NULL ,
        description TEXT NOT NULL 
    );
    `

	userTable := `
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY ASC ,
        email TEXT UNIQUE, 
		password TEXT,
		name TEXT,
		mfaSharedSecret TEXT,
		mfaEnrolled BOOLEAN DEFAULT 0,
		mfaBackupCodes TEXT
    );
    `

	jsonWebTokensTable := `
    CREATE TABLE IF NOT EXISTS jsonwebtoken(
    	id INTEGER PRIMARY KEY ASC,
    	keyId TEXT UNIQUE,
    	user INTEGER,
    	issuedAt DATE,
    	publicKey TEXT,
    	payload TEXT
    );
    `

	settingsTable := `
    CREATE TABLE IF NOT EXISTS settings(
        key TEXT PRIMARY KEY NOT NULL ,
        value TEXT NOT NULL 
    );
    `

	pushSubscriptionsTable := `
    CREATE TABLE IF NOT EXISTS pushSubscriptions(
        id INTEGER PRIMARY KEY ASC,
        keyId TEXT UNIQUE,
        user INTEGER,
        issuedAt DATE,
        publicKey TEXT,
        payload TEXT,
        subscription TEXT UNIQUE
    );
    `

	rulesTable := `
    CREATE TABLE IF NOT EXISTS rules(
        id INTEGER PRIMARY KEY ,
        description TEXT NOT NULL 
    );
    `
	_, err := s.db.Exec(thingsTable)
	_, err = s.db.Exec(userTable)
	_, err = s.db.Exec(jsonWebTokensTable)
	_, err = s.db.Exec(settingsTable)
	_, err = s.db.Exec(pushSubscriptionsTable)
	_, err = s.db.Exec(rulesTable)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) updateValue(k, v, table string) (err error) {
	query := fmt.Sprintf(`update %s set value=@value where key=@key`, table)
	_, err = s.db.Exec(query, sql.Named("value", v), sql.Named("key", k))
	return
}

func (s *Storage) queryValue(k string, table string) (value string, err error) {
	query := fmt.Sprintf("SELECT value FROM %s where key = @key", table)
	err = s.db.QueryRow(query, sql.Named("key", k)).Scan(&value)
	return value, err
}

func (s *Storage) deleteValue(key string, table string) error {
	stmt, err := s.db.Prepare(fmt.Sprintf(`delete from %s where key = ?`, table))
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) setValue(key string, value, table string) error {
	query := fmt.Sprintf("INSERT INTO %s(key, value) values(?,?)", table)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}(stmt)
	res, ee := stmt.Exec(key, value)
	if ee != nil {
		s.logger.Errorf("insert %s: key:%s err:", table, key, err.Error())
		return ee
	}
	_, eee := res.LastInsertId()
	if eee != nil {
		s.logger.Errorf("insert %s: key:%s err:", table, key, err.Error())
		return eee
	}
	s.logger.Debugf("insert %s,id: %s \t\n", table)
	return nil
}
