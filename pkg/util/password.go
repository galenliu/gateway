package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) {
	bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func GenerateEccKey() (publicKeyStr, privateKeyStr string, err error) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	priBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return
	}
	privateBlock := pem.Block{
		Type:  "ECD PRIVATE KEY",
		Bytes: priBytes,
	}
	bufferPrivate := new(bytes.Buffer)
	err = pem.Encode(bufferPrivate, &privateBlock)
	if err != nil {
		return
	}
	privateKeyStr = bufferPrivate.String()

	//生成公钥
	publicKey := privateKey.PublicKey
	bufferPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return
	}
	//pem
	publicBlock := pem.Block{
		Type:  "ecdsa public key",
		Bytes: bufferPublicKey,
	}
	bufferPublic := new(bytes.Buffer)
	err = pem.Encode(bufferPublic, &publicBlock)
	if err != nil {
		return "", "", err
	}
	publicKeyStr = bufferPublic.String()
	return
}