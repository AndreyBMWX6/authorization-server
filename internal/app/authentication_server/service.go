package authentication_server

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"oauth2/internal/app/scratch"
	"oauth2/internal/pkg/domain"
	"oauth2/internal/pkg/repositories"
	authenticateUser "oauth2/internal/pkg/usecases/authenticate_user"
	registerUser "oauth2/internal/pkg/usecases/register_user"
	desc "oauth2/pkg/api/authentication_server"
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
