package database

import (
	"database/sql"
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/things"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const dbFileName = "database.sqlite3"

var instance *store

type store struct {
	file string
	*sql.DB
}

type Store interface {
	things.Store
}

func (s *store) reset() {
	_, err := os.Stat(s.file)
	if err == nil {
		_ = os.Remove(s.file)
	}
}

func (s *store) createTable() error {

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
    CREATE TABLE IF NOT EXISTS jsonwebtokens(
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
	_, err := s.Exec(thingsTable)
	_, err = s.Exec(userTable)
	_, err = s.Exec(jsonWebTokensTable)
	_, err = s.Exec(settingsTable)
	_, err = s.Exec(pushSubscriptionsTable)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) updateValue(k, v string) (err error) {
	_, err = s.Exec(`update data set value=@value where key=@key`, sql.Named("value", v), sql.Named("key", k))
	return
}

func (s *store) queryValue(k string) (value string, err error) {
	err = s.QueryRow("SELECT value FROM data where key = @key", sql.Named("key", k)).Scan(&value)
	logging.Info(k, value)
	return value, err
}

func (s *store) deleteValue(key string) error {
	stmt, err := s.Prepare(`delete from data where key = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}
	return nil
}

func GetSetting(key string) (value string, err error) {

	err = instance.QueryRow("SELECT value FROM settings where key = @key", sql.Named("key", key), sql.Named("key", key)).Scan(&value)
	return value, err
}

func SetSetting(key, value string) error {

	logging.Info("set setting key:%s value:%s \t\n", key, value)
	_, err := GetSetting(key)
	if err == nil {
		_, e := instance.Exec(`update settings set value=@value where key=@key`, sql.Named("value", value), sql.Named("key", key))
		return e
	}
	stmt, err := instance.Prepare("INSERT INTO settings(key, value) values(?,?)")
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

func CreateThing(id, description string) error {
	if id == "" && description == "" {
		return fmt.Errorf("description is emtry")
	}
	stmt, err := instance.Prepare("INSERT INTO things(id, description) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, ee := stmt.Exec(id, description)
	if ee != nil {
		return ee
	}
	_, eee := res.LastInsertId()
	if eee != nil {
		return eee
	}
	fmt.Printf("insert data,id:%s , value: %s \t\n", id, description)
	return nil
}

func GetThings() map[string]string {
	var things = make(map[string]string)
	rows, err := instance.Query("SELECT id, description FROM things")
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id string
		var description string
		err = rows.Scan(&id, &description)
		if err == nil {
			things[id] = description
		}
	}
	return things
}

func RemoveThing(id string) error {
	stmt, err := instance.Prepare(`delete from things where id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil

}

func UpdateThing(id string, description string) (err error) {
	_, err = instance.Exec(`update things set id=@id where description=@description`, sql.Named("id", id), sql.Named("description", description))
	return
}

func GetUsersCount() map[string]string {
	var users = make(map[string]string)
	rows, err := instance.Query("SELECT id, description FROM users")
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id string
		var description string
		err = rows.Scan(&id, &description)
		if err == nil {
			users[id] = description
		}
	}
	return users
}

func CreateUser(email, hash, name string) (int64, error) {
	if email == "" {
		return 0, fmt.Errorf("email is emtry")
	}
	if hash == "" {
		return 0, fmt.Errorf("password hash error")
	}
	stmt, err := instance.Prepare("INSERT INTO users(email,password ,name ,mfaSharedSecret ,mfaEnrolled ,mfaBackupCodes) values(?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(stmt)
	res, ee := stmt.Exec(email, hash, name, "", false, "")
	if ee != nil {
		return 0, ee
	}
	id, eee := res.LastInsertId()
	if eee != nil {
		return id, eee
	}
	return id, nil
}
