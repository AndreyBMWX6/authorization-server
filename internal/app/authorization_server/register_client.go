package authorization_server

import (
	"context"
	"errors"

	desc "authorization-server/pkg/api/authorization_server"
)

// RegisterClient - registers client in authorization server
func (i *Implementation) RegisterClient(ctx context.Context, req *desc.RegisterClientRequest) (*desc.RegisterClientResponse, error) {
	return nil, errors.New("unimplemented")
}
