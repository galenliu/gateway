package db

import (
	"database/sql"
	"github.com/galenliu/gateway/pkg/logging"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

const dbFileName = "database.sqlite3"

type Store struct {
	file   string
	db     *sql.DB
	logger logging.Logger
}

func NewStore(filePath string, reset bool, log logging.Logger) (*Store, error) {
	s := &Store{}
	s.logger = log
	f := path.Join(filePath, dbFileName)
	s.file = f
	if reset {
		s.reset()
	}
	d, e := sql.Open("sqlite3", filePath)
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

func (s *Store) reset() {
	_, err := os.Stat(s.file)
	if err == nil {
		_ = os.Remove(s.file)
	}
}

func (s *Store) createTable() error {

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
	_, err := s.db.Exec(thingsTable)
	_, err = s.db.Exec(userTable)
	_, err = s.db.Exec(jsonWebTokensTable)
	_, err = s.db.Exec(settingsTable)
	_, err = s.db.Exec(pushSubscriptionsTable)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) updateValue(k, v string) (err error) {
	_, err = s.db.Exec(`update data set value=@value where key=@key`, sql.Named("value", v), sql.Named("key", k))
	return
}

func (s *Store) queryValue(k string) (value string, err error) {
	err = s.db.QueryRow("SELECT value FROM data where key = @key", sql.Named("key", k)).Scan(&value)
	return value, err
}

func (s *Store) deleteValue(key string) error {
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
