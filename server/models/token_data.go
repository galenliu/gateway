package models

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"time"
)

type TokenData struct {
	int64
	User      int64     `json:"user"`
	IssuedAt  time.Time `json:"issuedAt"`
	PublicKey string    `json:"publicKey"`
	KeyId     string    `json:"keyId"`
	Payload
}

func NewTokenData() *TokenData {
	return nil
}

func (t TokenData) privateKey() *ecdsa.PublicKey {
	block, _ := pem.Decode([]byte(t.PublicKey))
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}
	pr, ok := key.(ecdsa.PublicKey)
	if !ok {
		return nil
	}
	return &pr
}
