package authentication_server

import (
	"context"
	"net/http"

	"authorization-server/internal/app/scratch"
	"authorization-server/internal/pkg/domain"
	"authorization-server/internal/pkg/repositories"
	authenticateUser "authorization-server/internal/pkg/usecases/authenticate_user"
	registerUser "authorization-server/internal/pkg/usecases/register_user"
	desc "authorization-server/pkg/api/authentication_server"
	"github.com/jmoiron/sqlx"
)

type Implementation struct {
	desc.UnimplementedAuthenticationServerServer
	registerUserUseCase     RegisterUserUseCase
	authenticateUserUseCase AuthenticateUserUseCase
	fileServer              http.Handler
}

func NewAuthenticationServer(db *sqlx.DB, fileServer http.Handler) *Implementation {
	return &Implementation{
		registerUserUseCase:     registerUser.New(repositories.NewUsersRepository(db)),
		authenticateUserUseCase: authenticateUser.New(repositories.NewUsersRepository(db)),
		fileServer:              fileServer,
	}
}

func (i *Implementation) GetDescription() scratch.ServiceDesc {
	return desc.NewAuthorizationServerDesc(i)
}

type RegisterUserUseCase interface {
	Register(ctx context.Context, user *domain.User) error
}

type AuthenticateUserUseCase interface {
	Authenticate(ctx context.Context, user *domain.User) (string, error)
}
