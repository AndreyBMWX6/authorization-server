package claims

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	loginKey       = "name"
	expiredTimeKey = "exp"
)

var (
	ErrNotFound     = errors.New("claim not found")
	ErrWrongType    = errors.New("wrong claim value type")
	ErrParsingValue = errors.New("parsing claim value failed")
)

func GetLogin(claims jwt.MapClaims) (string, error) {
	loginClaim, ok := claims[loginKey]
	if !ok {
		return "", ErrNotFound
	}
	login, ok := loginClaim.(string)
	if !ok {
		return "", ErrWrongType
	}

	return login, nil
}

func GetExpiredTime(claims jwt.MapClaims) (time.Time, error) {
	expiredTimeClaim, ok := claims[expiredTimeKey]
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
