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

func (i *Implementation) RegisterUser(ctx context.Context, req *desc.RegisterUserRequest) (*desc.RegisterUserResponse, error) {
	if err := validateRegisterUserRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &domain.User{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}

	err := i.registerUserUseCase.Register(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "register user")
	}

	return &desc.RegisterUserResponse{}, nil
}

func validateRegisterUserRequest(ctx context.Context, req *desc.RegisterUserRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.Login, validation.Required),
		//todo: pass politics
		validation.Field(&req.Password, validation.Required),
	)
}
