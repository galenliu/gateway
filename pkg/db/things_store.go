package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
)

func (s *Storage) CreateThing(id string, thing *container.Thing) error {
	bytes, err := json.Marshal(thing)
	if err != nil {
		return err
	}
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
	return nil
}

func (s *Storage) GetThings() map[string]*container.Thing {
	things := make(map[string]*container.Thing)
	rows, err := s.db.Query("SELECT id, description FROM things")
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
		var thing container.Thing
		err := json.Unmarshal([]byte(description), &thing)
		if err != nil {
			continue
		}
		things[id] = &thing
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

func (s *Storage) UpdateThing(id string, thing *container.Thing) (err error) {
	bytes, err := json.Marshal(thing)
	if err != nil {
		return fmt.Errorf("thing marshal error: %s", err.Error())
	}
	_, err = s.db.Exec(`update things set id=@id where description=@description`, sql.Named("id", id), sql.Named("description", bytes))
	return
}
