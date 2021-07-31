package db

import (
	"database/sql"
	"fmt"
	"github.com/galenliu/gateway/server/models"
)



func (s *store) SaveThing(t *models.Thing) error {
	panic("implement me")
}


func (s *store)  CreateThing(id, description string) error {
	if id == "" && description == "" {
		return fmt.Errorf("description is emtry")
	}
	stmt, err := s.Prepare("INSERT INTO things(id, description) values(?,?)")
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

func (s *store)  GetThings() map[string]string {
	var things = make(map[string]string)
	rows, err := s.Query("SELECT id, description FROM things")
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

func (s *store)  RemoveThing(id string) error {
	stmt, err := s.Prepare(`delete from things where id = ?`)
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

func (s *store) UpdateThing(id string, description string) (err error) {
	_, err = s.Exec(`update things set id=@id where description=@description`, sql.Named("id", id), sql.Named("description", description))
	return
}


