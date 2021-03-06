package models

import (
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
)

type Users struct {
	users  []*User
	logger logging.Logger
}

func NewUsersModel(logger logging.Logger) *Users {
	users := &Users{}
	users.logger = logger
	return users
}

func (u *Users) GetUsersCount() int {
	if u.users != nil {
		return 0
	}
	return len(u.getUsersCount())
}

func (u *Users) GetUser(email string) *User {
	return nil
}

func (u *Users) CreateUser(email, password, name string) (error, string) {
	user := NewUser(email, password, name)
	err := user.Save()
	if err != nil {
		return err, ""
	}
	u.users = append(u.users, user)
	return nil, ""
}

func (u *Users) getUsersCount() []*User {
	database.GetUsersCount()
	return nil
}
