package authorization_server

import (
	"context"

	"github.com/pkg/errors"

	desc "authorization-server/pkg/api/authorization_server"
)

func (i *Implementation) GetAuthorizationCode(ctx context.Context, req *desc.GetAuthorizationCodeRequest) (*desc.GetAuthorizationCodeResponse, error) {
	return nil, errors.New("unimplemented")
}
