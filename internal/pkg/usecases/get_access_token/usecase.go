package get_access_token

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"oauth2/internal/pkg/domain"
	"oauth2/internal/pkg/jwt"
	"oauth2/internal/pkg/passwords"
)

type AuthorizationCodesRepository interface {
	GetCode(ctx context.Context, code string) (*domain.AuthorizationCode, error)
	MatchCodeUsed(ctx context.Context, code string) error
}

type ClientsRepository interface {
	GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error)
}

type TokensRepository interface {
	UpsertToken(ctx context.Context, token domain.Token) error
	GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
	DeleteTokensByCode(ctx context.Context, authorizationCode string) error
}

type UseCase struct {
	authorizationCodesRepo AuthorizationCodesRepository
	tokensRepo             TokensRepository
	clientsRepo            ClientsRepository
}

func New(
	authorizationCodesRepo AuthorizationCodesRepository,
	accessTokensRepo TokensRepository,
	clientsRepo ClientsRepository,
) *UseCase {
	return &UseCase{
		authorizationCodesRepo: authorizationCodesRepo,
		tokensRepo:             accessTokensRepo,
		clientsRepo:            clientsRepo,
	}
}

func (u *UseCase) GetTokenByAuthorizationCode(ctx context.Context, authorizationCode string, client domain.Client) (*domain.Token, error) {
	code, err := u.authorizationCodesRepo.GetCode(ctx, authorizationCode)
	if err != nil {
		return nil, errors.Wrap(err, "get authorization code meta")
	}

	if code.Used {
		// if code was already used revoking all issued tokens, return error
		err = u.tokensRepo.DeleteTokensByCode(ctx, authorizationCode)
		if err != nil {
			return nil, errors.Wrap(err, "delete tokens by code")
		}
		return nil, ErrUsedAuthorizationCode
	}

	err = u.authorizationCodesRepo.MatchCodeUsed(ctx, code.Code)
	if err != nil {
		return nil, errors.Wrap(err, "match code used")
	}

	if !time.Now().In(time.UTC).Before(code.ExpirationTime) {
		// todo: check if need to do smth
		return nil, ErrExpiredAuthorizationCode
	}

	if code.ClientID != client.ID {
		return nil, ErrAnotherClientAuthorizationCode
	}

	storedClient, err := u.clientsRepo.GetClient(ctx, client.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get client")
	}

	if storedClient.IsConfidential {
		if client.Secret == "" {
			return nil, ErrUnauthenticatedClient
		}

		match, err := passwords.MatchPasswords(storedClient.Secret, client.Secret)
		if err != nil {
			return nil, errors.Wrap(err, "match passwords")
		}
		if !match {
			return nil, ErrUnauthenticatedClient
		}
	}

	if code.RedirectURI != "" {
		if client.RedirectURI == "" {
			return nil, ErrNoRedirectURI
		}
		if client.RedirectURI != code.RedirectURI {
			return nil, ErrWrongRedirectURI
		}
	}

	accessToken, err := jwt.NewWithClaims(ctx, map[string]interface{}{
		"authorization_details": &domain.AuthorizationDetails{
			ClientID: client.ID,
			Scope:    code.Scope,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "build jwt token")
	}

	nonce := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "generate nonce for refresh token")
	}
	refreshToken, err := jwt.NewWithClaims(ctx, map[string]interface{}{
		"nonce": nonce,
	})

	//todo: move to config
	expiresIn := 3600 * time.Second
	token := domain.Token{
		AccessToken:       accessToken,
		AuthorizationCode: authorizationCode,
		Type:              domain.TokenTypeBearer,
		CreatedAt:         time.Now().In(time.UTC),
		ExpiresIn:         expiresIn,
		RefreshToken:      &refreshToken,
	}

	err = u.tokensRepo.UpsertToken(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "insert access token")
	}

	return &token, nil
}

func (u *UseCase) GetTokenByRefreshToken(ctx context.Context, refreshToken string, client domain.Client) (*domain.Token, error) {
	storedClient, err := u.clientsRepo.GetClient(ctx, client.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get client")
	}

	if storedClient.IsConfidential {
		if client.Secret == "" {
			return nil, ErrUnauthenticatedClient
		}

		match, err := passwords.MatchPasswords(storedClient.Secret, client.Secret)
		if err != nil {
			return nil, errors.Wrap(err, "match passwords")
		}
		if !match {
			return nil, ErrUnauthenticatedClient
		}
	}

	token, err := u.tokensRepo.GetTokenByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "get token by refresh token")
	}

	accessToken := token.AccessToken
	if !time.Now().In(time.UTC).Before(token.CreatedAt.Add(time.Second * token.ExpiresIn)) {
		accessToken, err = jwt.NewWithClaims(ctx, map[string]interface{}{
			"authorization_details": &domain.AuthorizationDetails{
				ClientID: client.ID,
				// todo: maybe rewrite to nil pointer scope
				Scope: lo.FromPtr(token.Scope),
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "build jwt token")
		}
	}

	// generate new refresh token
	nonce := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "generate nonce for refresh token")
	}
	refreshToken, err = jwt.NewWithClaims(ctx, map[string]interface{}{
		"nonce": nonce,
	})

	expiresIn := 3600 * time.Second
	token = &domain.Token{
		AccessToken:  accessToken,
		Type:         domain.TokenTypeBearer,
		CreatedAt:    time.Now().In(time.UTC),
		ExpiresIn:    expiresIn,
		RefreshToken: &refreshToken,
	}

	err = u.tokensRepo.UpsertToken(ctx, *token)
	if err != nil {
		return nil, errors.Wrap(err, "insert access token")
	}

	return token, nil
}
