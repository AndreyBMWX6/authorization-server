package authorization_server

import (
	"context"
	"errors"

	desc "authorization-server/pkg/api/authorization_server"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	return nil, errors.New("unimplemented")
}
