package authentication_server

import (
	"context"

	"authorization-server/internal/pkg/domain"
	desc "authorization-server/pkg/api/authentication_server"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) AuthenticateUser(ctx context.Context, req *desc.AuthenticateUserRequest) (*desc.AuthenticateUserResponse, error) {
	if err := validateAuthenticateUserRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &domain.User{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}

	jwt, err := i.authenticateUserUseCase.Authenticate(ctx, user)
	if err != nil {
		// todo: status code mapping
		return nil, errors.Wrap(err, "authenticate user")
	}

	return &desc.AuthenticateUserResponse{
		Jwt: jwt,
	}, nil
}

func validateAuthenticateUserRequest(ctx context.Context, req *desc.AuthenticateUserRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.Login, validation.Required),
		validation.Field(&req.Password, validation.Required),
	)
}
