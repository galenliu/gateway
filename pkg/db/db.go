package db

import (
	"database/sql"
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

func (s *Storage) updateValue(k, v string) (err error) {
	_, err = s.db.Exec(`update data set value=@value where key=@key`, sql.Named("value", v), sql.Named("key", k))
	return
}

func (s *Storage) queryValue(k string) (value string, err error) {
	err = s.db.QueryRow("SELECT value FROM data where key = @key", sql.Named("key", k)).Scan(&value)
	return value, err
}

func (s *Storage) deleteValue(key string) error {
	stmt, err := s.db.Prepare(`delete from data where key = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}
	return nil
}
