package middleware

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
)

type JWTMiddleware struct {
	logger logging.Logger
}

func NewJWTMiddleware(logger logging.Logger) *JWTMiddleware {
	j := &JWTMiddleware{}
	j.logger = logger
	return j
}

func (m JWTMiddleware) Handler(c *fiber.Ctx) error {
	jwtware.New(jwtware.Config{
		Filter:         nil,
		SuccessHandler: nil,
		ErrorHandler:   nil,
		SigningKey:     nil,
		SigningKeys:    nil,
		SigningMethod:  "",
		ContextKey:     "",
		Claims:         nil,
		jwt.Keyfunc(),
	})
	return nil
}
