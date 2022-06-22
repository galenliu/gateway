package middleware

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/galenliu/gateway/api/models"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"strings"
)

type JWTMiddleware struct {
	logger  log.Logger
	storage models.JsonwebtokenStore
}

func NewJWTMiddleware(storage models.JsonwebtokenStore) *JWTMiddleware {
	j := &JWTMiddleware{}
	j.storage = storage

	return j
}

func (m JWTMiddleware) Auth(c *fiber.Ctx) error {
	token, err := extractJWTHeader(c)
	var claimsMap models.Claims
	withClaims, err := jwt.ParseWithClaims(token, &claimsMap, func(t *jwt.Token) (any, error) {
		clm, ok := t.Claims.(*models.Claims)
		if !ok {
			return nil, errors.New("token invalid")
		}
		dbTokenData, err := m.storage.GetJSONWebTokenByKeyId(clm.KeyId)
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(dbTokenData.PublicKey)
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key invalid")
		}
		return publicKey, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("token invalid")
	}

	if withClaims.Valid {
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}

func extractJWTHeader(c *fiber.Ctx) (string, error) {
	var authScheme = "Bearer"
	auth := c.Get(fiber.HeaderAuthorization)
	l := len(authScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
		return auth[l+1:], nil
	}
	return "", errors.New("missing or malformed JWT")
}
