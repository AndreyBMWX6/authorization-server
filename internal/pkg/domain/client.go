package domain

import (
	"github.com/google/uuid"
)

type Client struct {
	ID          uuid.UUID
	Name        string
	URL         string
	RedirectURI string
	Secret      string
}
