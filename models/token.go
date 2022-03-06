package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ClientReadableToken struct {
	ExpiresAt int64    `json:"expires_at"`
	Roles     []string `json:"roles"`
}
type TokenClaims struct {
	jwt.RegisteredClaims
}

func MintToken(userid string, expires time.Time) (string, error) {
	claims := TokenClaims{
		jwt.RegisteredClaims{
			Issuer:    "gin-api",
			Subject:   userid,
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": claims.Issuer,
		"sub": claims.Subject,
		"exp": claims.ExpiresAt,
		"iat": claims.IssuedAt,
	})

	return token.SignedString([]byte(os.Getenv("DND_JWT_PRIVATE_KEY")))

}
