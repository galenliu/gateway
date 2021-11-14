package db

import (
	"database/sql"
	"fmt"
	json "github.com/json-iterator/go"
)

func (s *Storage) CreateThing(id string, thing interface{}) error {
	bytes, _ := json.MarshalIndent(thing, "", "  ")
	stmt, err := s.db.Prepare("INSERT INTO things(id, description) values(?,?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			s.logger.Error("stmt close err: %s", err.Error())
		}
	}(stmt)
	res, ee := stmt.Exec(id, bytes)
	if ee != nil {
		return ee
	}
	_, eee := res.LastInsertId()
	if eee != nil {
		return eee
	}
	fmt.Printf("insert data,id:%s , value: %s \t\n", id, bytes)
	return nil
}

func (s *Storage) GetThings() (things map[string][]byte) {
	if things == nil {
		things = make(map[string][]byte)
	}
	rows, err := s.db.Query("SELECT id, description FROM things")
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id string
		var description string
		err = rows.Scan(&id, &description)
		if err == nil {
			continue
		}
		things[id] = []byte(description)
	}
	return things
}

func (s *Storage) RemoveThing(id string) error {
	stmt, err := s.db.Prepare(`delete from things where id = ?`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			s.logger.Error("stmt close err: %s", err.Error())
		}
	}(stmt)
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateThing(id string, thing interface{}) (err error) {
	bytes, _ := json.Marshal(thing)
	_, err = s.db.Exec(`update things set id=@id where description=@description`, sql.Named("id", id), sql.Named("description", bytes))
	return
}
