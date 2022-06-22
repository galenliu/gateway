package models

import (
	"github.com/galenliu/gateway/pkg/db"
)

type UsersStore interface {
	CreateUser(user *db.User) (int64, error)
	DeleteUser(id int64) error
	GetUsers() []*db.User
	UpdateUser(user *db.User) error
}

type Users struct {
	store UsersStore
}

func NewUsersModel(store UsersStore) *Users {
	users := &Users{}
	users.store = store
	return users
}

func (u *Users) GetUsersCount() int {
	return len(u.store.GetUsers())
}

func (u *Users) GetUser(email string) *User {
	users := u.store.GetUsers()
	for _, user := range users {
		if user.Email == email {
			return &User{
				ID:              user.ID,
				Name:            user.Name,
				Email:           user.Email,
				Hash:            user.Hash,
				MfaSharedSecret: user.MfaSharedSecret,
				MfaEnrolled:     user.MfaEnrolled,
				MfaBackupCodes:  user.MfaBackupCodes,
			}
		}
	}
	return nil
}

func (u *Users) CreateUser(email, password, name string) (int64, error) {
	user := newUser(email, password, name)
	return u.store.CreateUser(&db.User{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Hash:            user.Hash,
		MfaSharedSecret: user.MfaSharedSecret,
		MfaEnrolled:     user.MfaEnrolled,
		MfaBackupCodes:  user.MfaBackupCodes,
	})
}
