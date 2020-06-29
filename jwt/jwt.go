package jwt

import (
	"errors"
	"time"

	jwtLib "github.com/dgrijalva/jwt-go"
)

// ErrMissingSecret error
var ErrMissingSecret = errors.New("missing secret key")

// Token wrapper
type Token struct {
	Secret  []byte
	Content string
}

// Claims for token
type Claims struct {
	Expiration time.Duration
	Fields     map[string]interface{}
	jwtLib.StandardClaims
}

// Generate JWT token for given claims
func (token *Token) Generate(claims *Claims) error {
	if token.Secret == nil || len(token.Secret) == 0 {
		return ErrMissingSecret
	}

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

	return claims, token.valid(tokenLib, claims)
}

// parse is helper function to parse jwt
func (token *Token) parse(tokenStr string, claims *Claims) (*jwtLib.Token, error) {
	tokenLib, err := jwtLib.ParseWithClaims(tokenStr, claims, func(t *jwtLib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, errors.New("failed to parse JWT token")
		}

		return token.Secret, nil
	})

	return tokenLib, err
}

// valid is helper function to validate jwt
func (token *Token) valid(tokenLib *jwtLib.Token, claims *Claims) bool {
	if tokenLib != nil {
		_, claimsOk := tokenLib.Claims.(jwtLib.Claims)

		return claimsOk && tokenLib.Valid && claims.ExpiresAt > time.Now().Unix()
	}

	return false
}
