package authorization_server

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"oauth2/internal/app/scratch"
	"oauth2/internal/pkg/domain"
	"oauth2/internal/pkg/repositories"
	getAccessToken "oauth2/internal/pkg/usecases/get_access_token"
	getAuthorizationCode "oauth2/internal/pkg/usecases/get_authorization_code"
	getUser "oauth2/internal/pkg/usecases/get_user"
	registerClient "oauth2/internal/pkg/usecases/register_client"
	desc "oauth2/pkg/api/authorization_server"
)

type Implementation struct {
	desc.UnimplementedAuthorizationServerServer
	getUserUseCase        GetUserUseCase
	registerClientUseCase RegisterClientUseCase
	getAuthorizationCode  GetAuthorizationCodeUseCase
	getAccessToken        GetAccessTokenUseCase
	fileServer            http.Handler
}

func NewAuthorizationServer(db *sqlx.DB, fileServer http.Handler) *Implementation {
	return &Implementation{
		getUserUseCase: getUser.New(
			repositories.NewUsersRepository(db),
		),
		registerClientUseCase: registerClient.New(
			repositories.NewClientsRepository(db),
		),
		getAuthorizationCode: getAuthorizationCode.New(
			repositories.NewClientsRepository(db),
			repositories.NewAuthorizationCodesRepository(db),
		),
		getAccessToken: getAccessToken.New(
			repositories.NewAuthorizationCodesRepository(db),
			repositories.NewTokensRepository(db),
			repositories.NewClientsRepository(db),
		),
		fileServer: fileServer,
	}
}

func (i *Implementation) GetDescription() scratch.ServiceDesc {
	return desc.NewAuthorizationServerDesc(i)
}

type GetUserUseCase interface {
	GetUser(ctx context.Context, login string) (*domain.User, error)
}

type RegisterClientUseCase interface {
	Register(ctx context.Context, client *domain.Client) (*domain.Client, error)
}

type GetAuthorizationCodeUseCase interface {
	GetCode(ctx context.Context, client *domain.Client, scope *string) (string, error)
}

type GetAccessTokenUseCase interface {
	GetTokenByAuthorizationCode(ctx context.Context, authorizationCode string, client domain.Client) (*domain.Token, error)
	GetTokenByRefreshToken(ctx context.Context, refreshToken string, client domain.Client) (*domain.Token, error)
}
