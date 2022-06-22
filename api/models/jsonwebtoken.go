package models

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/galenliu/gateway/pkg/db"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"time"
)

const RoleUserToken = "user_token"

type Payload struct {
	Role     string `json:"role"`
	Scope    string `json:"scope"`
	ClientId string `json:"clientId"`
}

type JsonwebtokenStore interface {
	CreateJSONWebToken(data *db.TokeDataStorage) error
	GetJSONWebTokenByKeyId(keyId string) (*db.TokeDataStorage, error)
}

type Jsonwebtoken struct {
	settingsModel *Settings
	Store         JsonwebtokenStore
}

func NewJsonwebtokenModel(settingsModel *Settings, store JsonwebtokenStore) *Jsonwebtoken {
	m := &Jsonwebtoken{}

	m.Store = store
	m.settingsModel = settingsModel
	return m
}

func (j *Jsonwebtoken) IssueToken(userId int64) (string, error) {
	sig, token, err := j.crateUser(userId, Payload{Role: RoleUserToken})
	if err != nil {
		return "", err
	}
	p, err := json.Marshal(token.Payload)
	err = j.Store.CreateJSONWebToken(&db.TokeDataStorage{
		KeyId:     token.KeyId,
		User:      token.User,
		IssuedAt:  token.IssuedAt,
		PublicKey: token.PublicKey,
		PayLoad:   p,
	})
	if err != nil {
		return "", fmt.Errorf("Issue token err : %s", err.Error())

	}
	return sig, nil
}

type Claims struct {
	UserId int64  `json:"userId"`
	KeyId  string `json:"KeyId"`
	Payload
	jwt.StandardClaims
}

func (j *Jsonwebtoken) crateUser(userId int64, payload Payload) (string, *TokenData, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return "", nil, err
	}
	keyId := uuid.NewV4().String()
	issuer := j.settingsModel.GetTunnelInfo()
	claims := &Claims{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{},
	}
	claims.KeyId = keyId
	claims.UserId = userId
	claims.Payload = payload
	if issuer != "" {
		claims.Issuer = issuer
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	sig, err := token.SignedString(privateKey)
	if err != nil {
		return "", nil, err
	}
	bufferPublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", nil, err
	}
	block := pem.Block{
		Type:  "ecdsa public key",
		Bytes: bufferPublicKey,
	}
	bufferPrivate := new(bytes.Buffer)
	err = pem.Encode(bufferPrivate, &block)
	if err != nil {
		return "", nil, err
	}
	tokenData := &TokenData{
		User:      userId,
		IssuedAt:  time.Now(),
		PublicKey: bufferPrivate.Bytes(),
		KeyId:     keyId,
		Payload:   payload,
	}
	return sig, tokenData, nil
}
