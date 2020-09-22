package util

import (
	"time"
)

var jwtKey = []byte("smart_assistant")

type Claims struct {
	UserId string
	jwt.StandardClaims
}

//生成一个Token
func ReleaseToken() (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: "id",
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: expirationTime.Unix(),
			Id:        "",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "smart_home",
			NotBefore: 0,
			Subject:   "user_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//
func ParseToken(tokenStr string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
