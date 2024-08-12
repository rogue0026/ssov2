package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrSigningKeyNotDefined = errors.New("signing key is not defined")
)

type tokenClaims struct {
	login string
	jwt.RegisteredClaims
}

func Generate(login string, ttl time.Duration) (string, error) {
	const fn = "internal.token.Generate"
	c := tokenClaims{
		login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	key := os.Getenv("KEY")
	if len(key) == 0 {
		return "", ErrSigningKeyNotDefined
	}
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return ss, nil
}
