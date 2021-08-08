package middleware

import (
	"github.com/galenliu/gateway/server/models"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	json "github.com/json-iterator/go"
	"strings"
)

func NewJWTWare(store models.JsonwebtokenStore) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		tokenString := ExtractJWTHeader(ctx)
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Token error")
		}
		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Token error")
		}

		// parse Claims
		var claimBytes []byte

		claimBytes, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Token error")
		}
		kid := json.Get(claimBytes, "keyId").ToString()
		token := store.GetJSONWebTokenByKeyId(kid)
		if token == nil {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Token error")
		}
		return jwtware.New(jwtware.Config{

			SigningKey: token.PublicKey,
		})(ctx)
	}

}

const AuthType = "Bearer"

func ExtractJWTHeader(c *fiber.Ctx) string {
	authorization := c.Get(fiber.HeaderAuthorization)
	if authorization == "" {
		return ""
	}
	l := strings.Split(authorization, " ")
	if l[0] != AuthType {
		return ""
	}
	var sig = ""
	if sig = l[1]; sig == "" {
		return ""
	}
	return sig
}
