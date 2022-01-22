package proxy

import (
	"database/sql"
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	json "github.com/json-iterator/go"
	_ "github.com/mattn/go-sqlite3"
)

func LoadConfig(path string, packageName string) (string, error) {
	db, err := sql.Open("sqlite3", path)
	defer db.Close()
	if err != nil {
		return "", err
	}
	var key = fmt.Sprintf("addons.config.%s", packageName)
	var value string
	err = db.QueryRow(`SELECT value FROM settings where key = @key`, sql.Named("key", key)).Scan(&value)
	return value, err
}

func SaveConfig(path string, packageName string, value string) error {
	db, err := sql.Open("sqlite3", path)
	defer db.Close()
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(`INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)`)
	defer stmt.Close()
	if err != nil {
		return err
	}
	var key = fmt.Sprintf("addons.config.%s", packageName)
	_, err = stmt.Exec(key, value)
	if err != nil {
		return err
	}
	return nil
}

func LoadThings(path string) map[string]things.Thing {
	db, err := sql.Open("sqlite3", path)
	defer db.Close()
	if err != nil {
		return nil
	}
	ts := make(map[string]things.Thing)
	rows, err := db.Query("SELECT id, description FROM things")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}(rows)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id string
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			continue
		}
		var thing things.Thing
		err := json.Unmarshal([]byte(description), &thing)
		if err != nil {
			continue
		}
		ts[id] = thing
	}
	return ts
}
