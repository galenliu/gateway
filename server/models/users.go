package models

import "github.com/galenliu/gateway/pkg/database"

type Users struct {
	users []*User
}

func NewUsers() *Users {
	users := &Users{}
	return users
}

func (u *Users) GetUsersCount() []*User {
	if u.users != nil {
		return u.users
	}
	return u.getUsersCount()
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
