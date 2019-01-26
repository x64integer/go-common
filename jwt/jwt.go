package jwt

import (
	"errors"
	"time"

	"github.com/x64integer/go-common/util"

	jwtLib "github.com/dgrijalva/jwt-go"
)

var (
	// LoginExp time
	LoginExp = 24 * time.Hour
	// claimsExp for jwt claims
	claimsExp = time.Now().Add(LoginExp).Unix()
	// secret for jwt
	secret = []byte(util.Env("JWT_LOGIN", util.RandomStr(64)))
)

// Generate JWT token for given claims
func Generate(claimOpts map[string]string) (string, error) {
	token := jwtLib.New(jwtLib.SigningMethodHS256)
	claims := token.Claims.(jwtLib.MapClaims)

	for key, value := range claimOpts {
		claims[key] = value
	}

	claims["exp"] = claimsExp

	tokenStr, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ValidateAndExtract will check if given JWT token is valid and return claims
func ValidateAndExtract(tokenStr string) (jwtLib.MapClaims, bool) {
	if token, _ := jwtLib.Parse(tokenStr, func(t *jwtLib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, errors.New("failed to parse JWT token")
		}

		return secret, nil
	}); token != nil {
		if claims, claimsOk := token.Claims.(jwtLib.MapClaims); claimsOk && token.Valid {
			return claims, true
		}
	}

	return nil, false
}
