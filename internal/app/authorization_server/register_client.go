package authorization_server

import (
	"context"

	"authorization-server/internal/pkg/domain"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "authorization-server/pkg/api/authorization_server"
)

// RegisterClient - registers client in authorization server
func (i *Implementation) RegisterClient(ctx context.Context, req *desc.RegisterClientRequest) (*desc.RegisterClientResponse, error) {
	if err := validateRegisterClientRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	client := mapClientToDomain(req)
	//todo: implement flow for public clients
	client.IsConfidential = true
	client, err := i.registerClientUseCase.Register(ctx, client)
	if err != nil {
		return nil, errors.Wrap(err, "register client")
	}

	return &desc.RegisterClientResponse{
		ClientId:     client.ID.String(),
		ClientSecret: client.Secret,
	}, nil
}

func validateRegisterClientRequest(ctx context.Context, req *desc.RegisterClientRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	return validation.ValidateStructWithContext(
		ctx,
		req,
		validation.Field(&req.Name, validation.Required),
		//todo: проверить как работает с https
		//todo: return
		//validation.Field(&req.Url, is.URL),
		//todo: проверить, как работает
		//todo: return
		//validation.Field(&req.RedirectUri, is.RequestURI),
	)
}

func mapClientToDomain(req *desc.RegisterClientRequest) *domain.Client {
	return &domain.Client{
		Name:        req.GetName(),
		URL:         req.GetUrl(),
		RedirectURI: req.GetRedirectUri(),
	}
}
