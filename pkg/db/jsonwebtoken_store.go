package db

import (
	"database/sql"
	"time"
)

type TokeDataStorage struct {
	KeyId     string
	User      int64
	IssuedAt  time.Time
	PublicKey []byte
	PayLoad   []byte
}

func (s *Storage) CreateJSONWebToken(t *TokeDataStorage) error {
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
	_, err = stmt.Exec(t.KeyId, t.User, t.IssuedAt, t.PublicKey, t.PayLoad)
	return err
}

func (s *Storage) GetJSONWebTokenByKeyId(keyId string) *TokeDataStorage {
	panic("implement me")
}
