package database

import (
	"database/sql"
	"fmt"
	"gateway/log"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

const dbFileName = "database.sqlite3"

var database *sql.DB

var filePath string

func InitDB(dir string) error {
	var e error
	filePath = path.Join(dir, dbFileName)
	database, e = sql.Open("sqlite3", filePath)
	if e != nil {
		return e
	}
	ee := createTable(database)
	if ee != nil {
		return ee
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

func createTable(db *sql.DB) error {

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
	_, err := db.Exec(thingsTable)
	_, err = db.Exec(userTable)
	_, err = db.Exec(jsonWebTokensTable)
	_, err = db.Exec(settingsTable)
	_, err = db.Exec(pushSubscriptionsTable)
	if err != nil {
		return err
	}
	return nil
}

func UpdateValue(k, v string) (err error) {
	_, err = database.Exec(`update data set value=@value where key=@key`, sql.Named("value", v), sql.Named("key", k))
	return
}

func QueryValue(k string) (value string, err error) {
	err = database.QueryRow("SELECT value FROM data where key = @key", sql.Named("key", k)).Scan(&value)
	log.Info(k, value)
	return value, err
}

func DeleteValue(key string) error {
	stmt, err := database.Prepare(`delete from data where key = ?`)
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

	err = database.QueryRow("SELECT value FROM settings where key = @key", sql.Named("key", key), sql.Named("key", key)).Scan(&value)
	log.Info("get setting key:%v value:%v", key, value)
	return value, err
}

func SetSetting(key, value string) error {

	log.Info("set setting key:%v value:%v", key, value)
	_, err := GetSetting(key)
	if err == nil {
		_, e := database.Exec(`update settings set value=@value where key=@key`, sql.Named("value", value), sql.Named("key", key))
		return e
	}
	stmt, err := database.Prepare("INSERT INTO settings(key, value) values(?,?)")
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
	log.Debug("insert data,id:%d , value: %s", id, value)
	return nil
}

func CreateThing(id, description string) error {
	if id == "" && description == "" {
		return fmt.Errorf("description is emtry")
	}
	stmt, err := database.Prepare("INSERT INTO things(id, description) values(?,?)")
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
	fmt.Printf("insert data,id:%d , value: %s", id, description)
	return nil
}

func GetThings() map[string]string {
	var things map[string]string = make(map[string]string)
	rows, err := database.Query("SELECT id, description FROM things")
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
	stmt, err := database.Prepare(`delete from things where id = ?`)
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
	_, err = database.Exec(`update things set id=@id where description=@description`, sql.Named("id", id), sql.Named("description", description))
	return
}
