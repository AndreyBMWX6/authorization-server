package authorization_server

import (
	"context"
	"net/http"

	"authorization-server/internal/app/scratch"
	"authorization-server/internal/pkg/domain"
	"authorization-server/internal/pkg/repositories"
	getUser "authorization-server/internal/pkg/usecases/get_user"
	registerClient "authorization-server/internal/pkg/usecases/register_client"
	desc "authorization-server/pkg/api/authorization_server"
	"github.com/jmoiron/sqlx"
)

type Implementation struct {
	desc.UnimplementedAuthorizationServerServer
	getUserUseCase        GetUserUseCase
	registerClientUseCase RegisterClientUseCase
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
