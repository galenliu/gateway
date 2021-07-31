package util

import "golang.org/x/crypto/bcrypt"

func Hash(password string) {
	bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
