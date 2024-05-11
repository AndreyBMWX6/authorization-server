package authenticate_user

import (
	"github.com/pkg/errors"
)

var (
	ErrWrongPassword = errors.New("wrong password")
)
