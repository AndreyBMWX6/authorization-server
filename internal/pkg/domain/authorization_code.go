package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationCode struct {
	Code           string
	ClientID       uuid.UUID
	RedirectURI    string
	ExpirationTime time.Time
	Scope          string
	Used           bool
}
