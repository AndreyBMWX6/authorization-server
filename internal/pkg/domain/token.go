package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeBearer TokenType = "Bearer"
)

type Token struct {
	AccessToken       string
	AuthorizationCode string
	Type              TokenType
	CreatedAt         time.Time
	// the lifetime in seconds of the access token
	ExpiresIn    time.Duration
	RefreshToken *string
	Scope        *string
}

// AuthorizationDetails - used in jwt claims
type AuthorizationDetails struct {
	ClientID uuid.UUID
	Scope    string
}
