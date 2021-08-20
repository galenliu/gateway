package db

import (
	"database/sql"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
)

func (s *Storage) CreateJSONWebToken(t *models.TokenData) error {
	stmt, err := s.db.Prepare("INSERT INTO jsonwebtoken(keyId,user,issuedAt,publicKey,payload) values(?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}(stmt)
	p, _ := json.MarshalToString(t.Payload)
	_, err = stmt.Exec(t.KeyId, t.User, t.IssuedAt, t.PublicKey, p)
	return err
}

func (s *Storage) GetJSONWebTokenByKeyId(keyId string) *models.TokenData {
	panic("implement me")
}
