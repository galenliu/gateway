package db

import (
	"database/sql"
	"fmt"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
)

func (s *Storage) SaveThing(t *models.Thing) error {
	panic("implement me")
}

func (s *Storage) CreateThing(t *models.Thing) error {
	if t.ID == "" {
		return fmt.Errorf("description is emtry")
	}
	description, _ := json.MarshalToString(t)
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
	res, ee := stmt.Exec(t.GetID(), description)
	if ee != nil {
		return ee
	}
	_, eee := res.LastInsertId()
	if eee != nil {
		return eee
	}
	fmt.Printf("insert data,id:%s , value: %s \t\n", t.GetID(), description)
	return nil
}

func (s *Storage) GetThings() (things []*models.Thing) {
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
		t, er := models.NewThingFromString(description)
		if er != nil {
			s.logger.Error("thing err: %s", er.Error())
			continue
		}
		things = append(things, t)
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

func (s *Storage) UpdateThing(t *models.Thing) (err error) {
	d, _ := json.MarshalToString(t)
	_, err = s.db.Exec(`update things set id=@id where description=@description`, sql.Named("id", t.GetID()), sql.Named("description", d))
	return
}
