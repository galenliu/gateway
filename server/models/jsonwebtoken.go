package models

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Payload struct {
	Role     string `json:"role"`
	Scope    string `json:"scope"`
	ClientId string `json:"clientId"`
}

type TokenData struct {
	User      int64     `json:"user"`
	IssuedAt  time.Time `json:"issuedAt"`
	PublicKey string    `json:"publicKey"`
	KeyId     string    `json:"keyId"`
	Payload
}

type JsonwebtokenStore interface {
	CreateJSONWebToken(data *TokenData) error
}

type Jsonwebtoken struct {
	settingsModel *Settings
	store         JsonwebtokenStore
	logger        logging.Logger
}

func (j *Jsonwebtoken) IssueToken(user int64) string {
	sig, token, err := j.crateUser(user, Payload{Role: RoleUserToken})
	if err != nil {
		j.logger.Info("Issue token err : %s", err.Error())
	}
	err = j.store.CreateJSONWebToken(token)
	if err != nil {
		j.logger.Info("Issue token err : %s", err.Error())

	}
	return sig
}

func NewJsonwebtokenModel(settingsModel *Settings, store JsonwebtokenStore, logger logging.Logger) *Jsonwebtoken {
	m := &Jsonwebtoken{}
	m.logger = logger
	m.settingsModel = settingsModel
	return m
}

type Claims struct {
	jwt.StandardClaims
	Issuer string `json:"issuer"`
	KeyId  string `json:"keyId"`
	Payload
}

const RoleUserToken = "user_token"

func (j *Jsonwebtoken) crateUser(userId int64, payload Payload) (string, *TokenData, error) {
	publicKeyStr, privateKeyStr, err := util.GenerateEccKey()

	keyId := uuid.NewV4()
	issuer := j.settingsModel.GetTunnelInfo()

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{},
		Issuer:         issuer,
		KeyId:          keyId.String(),
	}
	claims.Payload = payload

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	t, err := token.SignedString(privateKeyStr)
	if err != nil {
		return "", nil, err
	}
	tokenData := &TokenData{
		User:      userId,
		IssuedAt:  time.Now(),
		PublicKey: publicKeyStr,
		KeyId:     keyId.String(),
		Payload:   payload,
	}
	return t, tokenData, nil
}
