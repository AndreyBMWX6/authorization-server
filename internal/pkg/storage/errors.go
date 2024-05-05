package storage

import (
	"github.com/pkg/errors"
)

var (
	ErrAlreadyExists = errors.New("entity already exist")
	ErrNotFound      = errors.New("entity not found")
)
