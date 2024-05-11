package secret

import (
	"context"
)

const (
	JWTSecretKey = "jwt_secret"
)

// todo: rewrite to vault storage
type secretProvider struct {
	m map[string]interface{}
}

var provider = secretProvider{
	m: make(map[string]interface{}),
}

func init() {
	// define secret values in this package in secrets.go file
	provider.m[JWTSecretKey] = jwtSecretValue
}

func GetValue(ctx context.Context, key string) (interface{}, error) {
	value, ok := provider.m[key]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}
