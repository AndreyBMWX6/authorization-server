package get_access_token

import (
	"github.com/pkg/errors"
)

var (
	ErrUsedAuthorizationCode          = errors.New("authorization code was already used")
	ErrExpiredAuthorizationCode       = errors.New("authorization code expired")
	ErrAnotherClientAuthorizationCode = errors.New("authorization code was issues to another client")

	ErrUnauthenticatedClient = errors.New("unauthenticated client")

	ErrNoRedirectURI    = errors.New("no redirect uri provided by client")
	ErrWrongRedirectURI = errors.New("wrong redirect provided by client")
)
