package authorization_server

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "authorization-server/pkg/api/authorization_server"
)

func (i *Implementation) GetAuthorizationCode(ctx context.Context, req *desc.GetAuthorizationCodeRequest) (*desc.GetAuthorizationCodeResponse, error) {
	if err := validateGetAuthorizationCodeRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return nil, errors.New("unimplemented")
}

func validateGetAuthorizationCodeRequest(ctx context.Context, req *desc.GetAuthorizationCodeRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.ResponseType, validation.Required),
		validation.Field(&req.ClientId, validation.Required, is.UUIDv4),
	)
}
