package util

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// LoginExp time
	LoginExp  = 24 * time.Hour
	claimsExp = time.Now().Add(LoginExp).Unix()
	secret    = []byte(Env("JWT_LOGIN", RandomStr(64)))
)

// Token structure
type Token struct{}

// Generate JWT token for given claims
func (t *Token) Generate(c map[string]string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range c {
		claims[k] = v
	}

	claims["exp"] = claimsExp

	tokenStr, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ValidAndExtract will check if given JWT token is valid and return claims
func (t *Token) ValidAndExtract(tokenStr string) (jwt.MapClaims, bool) {
	if token, _ := jwt.Parse(tokenStr, func(_t *jwt.Token) (interface{}, error) {
		if _, ok := _t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("failed to parse JWT token")
		}

		return secret, nil
	}); token != nil {
		if claims, claimsOk := token.Claims.(jwt.MapClaims); claimsOk && token.Valid {
			return claims, true
		}
	}

	return nil, false
}
