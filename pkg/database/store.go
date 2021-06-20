package database

import (
	"database/sql"
	"path"
)

func NewStore(filePath string, reset bool) (Store, error) {
	s := &store{}

	fn := path.Join(filePath, dbFileName)
	s.file = fn

	d, e := sql.Open("sql	ite3", filePath)
	if e != nil {

		return nil, e
	}
	s.DB = d
	ee := s.createTable()
	if ee != nil {
		return nil, ee
	}
	if reset {
		s.reset()
	}
	return s, nil
}
