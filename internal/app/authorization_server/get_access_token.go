package authorization_server

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oauth2/internal/pkg/domain"

	desc "oauth2/pkg/api/authorization_server"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	if err := validateGetAccessTokenRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	switch req.GetGrantType() {
	case desc.GrantType_authorization_code:
		return i.getAccessTokenByAuthorizationCode(ctx, req)
	case desc.GrantType_refresh_token:
		return i.getAccessTokenByRefreshToken(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, "unknown grant type")
	}
}

func (i *Implementation) getAccessTokenByAuthorizationCode(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	if err := validateGetAccessTokenByAuthorizationCodeRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, errors.Wrap(err, "parse clientID")
	}

	// fix encoding bug
	code := strings.Replace(req.GetCode(), " ", "+", -1)
	client := domain.Client{
		ID:          clientID,
		RedirectURI: req.GetRedirectUri(),
		Secret:      req.GetClientSecret(),
	}

	token, err := i.getAccessToken.GetTokenByAuthorizationCode(ctx, code, client)
	if err != nil {
		return nil, errors.Wrap(err, "get access token by authorization code")
	}

	return &desc.GetAccessTokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    string(token.Type),
		ExpiresIn:    int64(token.ExpiresIn.Seconds()),
		RefreshToken: token.RefreshToken,
		Scope:        token.Scope,
	}, nil
}

func (i *Implementation) getAccessTokenByRefreshToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	if err := validateGetAccessTokenByRefreshTokenRequest(ctx, req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, errors.Wrap(err, "parse clientID")
	}
	client := domain.Client{
		ID:     clientID,
		Secret: req.GetClientSecret(),
	}

	token, err := i.getAccessToken.GetTokenByRefreshToken(ctx, req.GetRefreshToken(), client)
	if err != nil {
		return nil, errors.Wrap(err, "get access token by authorization code")
	}

	return &desc.GetAccessTokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    string(token.Type),
		ExpiresIn:    int64(token.ExpiresIn.Seconds()),
		RefreshToken: token.RefreshToken,
		Scope:        token.Scope,
	}, nil
}

func validateGetAccessTokenRequest(ctx context.Context, req *desc.GetAccessTokenRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.GrantType, validation.Required),
	)
}

func validateGetAccessTokenByAuthorizationCodeRequest(ctx context.Context, req *desc.GetAccessTokenRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.Code, validation.Required),
	)
}

func validateGetAccessTokenByRefreshTokenRequest(ctx context.Context, req *desc.GetAccessTokenRequest) error {
	return validation.ValidateStructWithContext(ctx, req,
		validation.Field(&req.RefreshToken, validation.Required),
	)
}
