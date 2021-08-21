package models

import (
	"encoding/json"
	"github.com/galenliu/gateway/pkg/db"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
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

type Claims struct {
	jwt.StandardClaims
	KeyId string `json:"keyId"`
	Payload
}

type TokenData struct {
	int64
	User      int64     `json:"user"`
	IssuedAt  time.Time `json:"issuedAt"`
	PublicKey string    `json:"publicKey"`
	KeyId     string    `json:"keyId"`
	Payload
}

type JsonwebtokenStore interface {
	CreateJSONWebToken(data *db.TokeDataStorage) error
	GetJSONWebTokenByKeyId(keyId string) *db.TokeDataStorage
}

type Jsonwebtoken struct {
	settingsModel *Settings
	Store         JsonwebtokenStore
	logger        logging.Logger
}

func NewJsonwebtokenModel(settingsModel *Settings, store JsonwebtokenStore, logger logging.Logger) *Jsonwebtoken {
	m := &Jsonwebtoken{}
	m.logger = logger
	m.Store = store
	m.settingsModel = settingsModel
	return m
}

func (j *Jsonwebtoken) IssueToken(user int64) string {
	sig, token, _ := j.crateUser(user, Payload{Role: RoleUserToken})
	p, err := json.Marshal(token.Payload)
	err = j.Store.CreateJSONWebToken(&db.TokeDataStorage{
		KeyId:     token.KeyId,
		User:      token.User,
		IssuedAt:  token.IssuedAt,
		PublicKey: token.PublicKey,
		PayLoad:   p,
	})
	if err != nil {
		j.logger.Info("Issue token err : %s", err.Error())
		return ""
	}
	return sig
}

func (j *Jsonwebtoken) crateUser(userId int64, payload Payload) (string, *TokenData, error) {
	publicKeyStr, privateKeyStr, err := util.GenerateEccKey()
	if err != nil {
		return "", nil, err
	}
	keyId := uuid.NewV4()
	issuer := j.settingsModel.GetTunnelInfo()
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{},
		KeyId:          keyId.String(),
	}
	if issuer != "" {
		claims.Issuer = issuer
	}
	claims.Payload = payload
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	sig, err := token.SignedString(privateKeyStr)
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
	return sig, tokenData, nil
}
