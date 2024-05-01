package authorization_server

import (
	"context"

	desc "authorization-server/pkg/api/authorization_server"
)

// RegisterClient - registers client in authorization server
func (i *Implementation) RegisterClient(ctx context.Context, req *desc.RegisterClientRequest) (*desc.RegisterClientResponse, error) {
	return &desc.RegisterClientResponse{}, nil
}
