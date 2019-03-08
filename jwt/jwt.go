package jwt

import (
	"errors"
	"sync"
	"time"

	jwtLib "github.com/dgrijalva/jwt-go"
)

// Token wrapper
type Token struct {
	Secret  []byte
	Content string
	Lock    sync.Mutex
}

// Claims for token
type Claims struct {
	Expiration time.Duration
	Fields     interface{}
	jwtLib.StandardClaims
}

// Generate JWT token for given claims
func (token *Token) Generate(claims *Claims) error {
	token.Lock.Lock()
	defer token.Lock.Unlock()

	claims.StandardClaims.ExpiresAt = time.Now().Add(claims.Expiration).Unix()

	tokenLib := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)

	tokenStr, err := tokenLib.SignedString(token.Secret)
	if err != nil {
		return err
	}

	token.Content = tokenStr

	return nil
}

// ValidateAndExtract will check if given JWT token is valid and return claims
func (token *Token) ValidateAndExtract(tokenStr string) (*Claims, bool) {
	claims := &Claims{}

	tokenLib, err := token.parse(tokenStr, claims)
	if err != nil {
		return claims, false
	}

	return claims, token.valid(tokenLib)
}

// parse is helper function to parse jwt
func (token *Token) parse(tokenStr string, claims *Claims) (*jwtLib.Token, error) {
	tokenLib, err := jwtLib.ParseWithClaims(tokenStr, claims, func(t *jwtLib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, errors.New("failed to parse JWT token")
		}

		return token.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	return tokenLib, nil
}

// valid is helper function to validate jwt
func (token *Token) valid(tokenLib *jwtLib.Token) bool {
	if tokenLib != nil {
		_, claimsOk := tokenLib.Claims.(jwtLib.Claims)

		if claimsOk && tokenLib.Valid {
			return true
		}
	}

	return false
}
