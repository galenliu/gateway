package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID              int64
	Name            string
	Email           string
	Hash            string
	MfaSharedSecret string
	MfaEnrolled     bool
	MfaBackupCodes  string
}

func (s *Storage) GetUsers() []*User {
	var users []*User
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		s.logger.Error("DataBase GetUsers Err:", err.Error())
		return nil
	}
	for rows.Next() {
		var (
			id              int64
			email           string
			password        string
			name            string
			mfaSharedSecret string
			mfaEnrolled     bool
			mfaBackupCodes  string
		)
		err = rows.Scan(&id, &email, &password, &name, &mfaSharedSecret, &mfaEnrolled, &mfaBackupCodes)
		if err == nil {
			users = append(users, &User{
				ID:              id,
				Name:            name,
				Email:           email,
				Hash:            password,
				MfaSharedSecret: mfaSharedSecret,
				MfaEnrolled:     mfaEnrolled,
				MfaBackupCodes:  mfaBackupCodes,
			})
		}
	}
	return users
}

func (s *Storage) CreateUser(u *User) (int64, error) {
	if u.Email == "" {
		return 0, fmt.Errorf("email is emtry")
	}
	if u.Hash == "" {
		return 0, fmt.Errorf("password hash error")
	}
	stmt, err := s.db.Prepare("INSERT INTO users(email,password ,name ,mfaSharedSecret ,mfaEnrolled ,mfaBackupCodes) values(?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}(stmt)
	res, ee := stmt.Exec(u.Email, u.Hash, u.Name, u.MfaSharedSecret, u.MfaEnrolled, u.MfaBackupCodes)
	if ee != nil {
		return 0, ee
	}
	id, eee := res.LastInsertId()
	if eee != nil {
		return id, eee
	}
	return id, nil
}

func (s *Storage) DeleteUser(id int64) error {
	stmt, err := s.db.Prepare("delete from users where id=?")
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}(stmt)
	if err != nil {
		return err
	}
	_, e := stmt.Exec(id)
	if e != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateUser(u *User) error {
	stmt, err := s.db.Prepare("update users set password=? name=? mfaSharedSecret=? mfaEnrolled=? mfaBackupCodes=? where id=?")
	if err != nil {
		return err
	}
	_, e := stmt.Exec(u.Hash, u.Name, u.MfaSharedSecret, u.MfaEnrolled, u.MfaBackupCodes, u.ID)
	if e != nil {
		return e
	}
	return nil
}
