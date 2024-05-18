package jwt

import (
	"context"
	"time"

	"authorization-server/internal/config/secret"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	loginKey       = "name"
	expiredTimeKey = "exp"
)

var (
	ErrNilToken     = errors.New("token is nil")
	ErrNotFound     = errors.New("claim not found")
	ErrWrongType    = errors.New("wrong claim value type")
	ErrParsingValue = errors.New("parsing claim value failed")
)

func GetLogin(token *jwt.Token) (string, error) {
	c, err := getClaimsMap(token)
	if err != nil {
		return "", errors.Wrap(err, "get claims map")
	}

	loginClaim, ok := c[loginKey]
	if !ok {
		return "", ErrNotFound
	}
	login, ok := loginClaim.(string)
	if !ok {
		return "", ErrWrongType
	}

	return login, nil
}

func GetExpiredTime(token *jwt.Token) (time.Time, error) {
	c, err := getClaimsMap(token)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "get claims map")
	}

	expiredTimeClaim, ok := c[expiredTimeKey]
	if !ok {
		return time.Time{}, ErrNotFound
	}
	expiredTimeStr, ok := expiredTimeClaim.(string)
	if !ok {
		return time.Time{}, ErrWrongType
	}

	expiredTime, err := time.Parse(time.RFC3339, expiredTimeStr)
	if !ok {
		return time.Time{}, errors.Wrap(err, ErrParsingValue.Error())
	}

	return expiredTime, nil
}

func NewWithClaims(ctx context.Context, claimsMap map[string]interface{}) (string, error) {
	//todo: move expires in to config
	expiresIn := 3600 * time.Second
	expirationTime := time.Now().In(time.UTC).Add(expiresIn)

	// basic claims
	claims := jwt.MapClaims{
		"exp": expirationTime,
	}
	for key, val := range claimsMap {
		claims[key] = val
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecretVal, err := secret.GetValue(ctx, secret.JWTSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "get jwt secret")
	}
	jwtSecret := jwtSecretVal.(string)

	return token.SignedString([]byte(jwtSecret))
}

func getClaimsMap(token *jwt.Token) (jwt.MapClaims, error) {
	if token == nil {
		return nil, ErrNilToken
	}

	c, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cast claims to map claims")
	}

	return c, nil
}
