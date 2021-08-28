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

func (s *Storage) GetJSONWebTokenByKeyId(keyId string) (*TokeDataStorage, error) {
	row := s.db.QueryRow("SELECT user, issuedAt, publicKey, payload from jsonwebtoken where keyId =@keyId", sql.Named("keyId", keyId))
	t := TokeDataStorage{}
	t.KeyId = keyId
	if err := row.Scan(&t.User, &t.IssuedAt, &t.PublicKey, &t.PayLoad); err != nil {
		return nil, err
	}
	return &t, nil
}
