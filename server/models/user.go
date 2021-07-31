package models

import (
	"golang.org/x/crypto/bcrypt"
	"log"
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

func newUser(email, password, name string) *User {
	u := &User{}
	hash := hashAndSalt(password)
	u.Hash = hash
	u.Name = name
	u.Email = email
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
