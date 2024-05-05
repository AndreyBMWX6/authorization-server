package register_client

import (
	"context"

	"authorization-server/internal/pkg/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
)

type ClientsRepository interface {
	InsertClient(ctx context.Context, client domain.Client) error
}

type UseCase struct {
	clientsRepo ClientsRepository
}

func New(clientsRepo ClientsRepository) *UseCase {
	return &UseCase{
		clientsRepo: clientsRepo,
	}
}

func (u *UseCase) Register(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	client.ID = uuid.New().String()

	secret, err := password.Generate(64, 10, 10, false, true)
	if err != nil {
		return nil, errors.Wrap(err, "generate client secret")
	}
	client.Secret = secret

	err = u.clientsRepo.InsertClient(ctx, *client)
	if err != nil {
		return nil, errors.Wrap(err, "insert client")
	}

	return client, nil
}
