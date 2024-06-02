package authorization_server

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oauth2/internal/pkg/domain"

	desc "oauth2/pkg/api/authorization_server"
)

func (i *Implementation) GetAuthorizationCode(ctx context.Context, req *desc.GetAuthorizationCodeRequest) (*desc.GetAuthorizationCodeResponse, error) {
	if err := validateGetAuthorizationCodeRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	client := &domain.Client{
		ID:          clientID,
		Name:        "",
		URL:         "",
		RedirectURI: req.GetRedirectUri(),
		Secret:      "",
	}
	code, err := i.getAuthorizationCode.GetCode(ctx, client, req.Scope)
	if err != nil {
		//todo: map to status err
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.GetAuthorizationCodeResponse{
		Code:  code,
		State: req.State,
	}, nil
}

func validateGetAuthorizationCodeRequest(ctx context.Context, req *desc.GetAuthorizationCodeRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.ResponseType, validation.Required),
		validation.Field(&req.ClientId, validation.Required, is.UUIDv4),
	)
}
