package wot_models

import (
	"github.com/galenliu/gateway/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	ID    int64
	Name  string
	Email string
	Hash  string
}

func NewUser(email, password, name string) *User {
	u := &User{}
	hash := hashAndSalt(password)
	u.Hash = hash
	u.Name = name
	u.Email = email
	return nil
}

// Save user保存到数据库中
func (u *User) Save() error {
	id, err := db.CreateUser(u.Email, u.Hash, u.Name)
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
