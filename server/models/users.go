package models

import (
	"github.com/galenliu/gateway/pkg/logging"
)

type UsersStore interface {
	CreateUser(user *User) (int64, error)
	DeleteUser(id int64) error
	GetUsers() []*User
	UpdateUser(user *User) error
}

type Users struct {
	logger logging.Logger
	store  UsersStore
}

func NewUsersModel(store UsersStore, logger logging.Logger) *Users {
	users := &Users{}
	users.store = store
	users.logger = logger
	return users
}

func (u *Users) GetUsersCount() int {
	return len(u.store.GetUsers())
}

func (u *Users) GetUser(email string) *User {
	users := u.store.GetUsers()
	for _, user := range users {
		if user.Email == email {
			return user
		}
	}
	return nil
}

func (u *Users) CreateUser(email, password, name string) (int64, error) {
	user := newUser(email, password, name)
	return u.store.CreateUser(user)
}
