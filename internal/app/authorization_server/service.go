package authorization_server

import (
	"context"

	"authorization-server/internal/app/scratch"
	"authorization-server/internal/pkg/domain"
	"authorization-server/internal/pkg/repositories"
	"authorization-server/internal/pkg/storage"
	registerClient "authorization-server/internal/pkg/usecases/register_client"
	desc "authorization-server/pkg/api/authorization_server"
)

type Implementation struct {
	desc.UnimplementedAuthorizationServerServer
	registerClientUseCase RegisterClientUseCase
}

func NewAuthorizationServer(storage *storage.Storage) *Implementation {
	return &Implementation{
		registerClientUseCase: registerClient.New(
			repositories.NewClientsRepository(storage),
		),
	}
}

func (i *Implementation) GetDescription() scratch.ServiceDesc {
	return desc.NewAuthorizationServerDesc(i)
}

type RegisterClientUseCase interface {
	Register(ctx context.Context, client *domain.Client) (*domain.Client, error)
}
