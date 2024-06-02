package register_client

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"oauth2/internal/pkg/domain"
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
	client.ID = uuid.New()
	secret, err := password.Generate(64, 10, 10, false, true)
	if err != nil {
		return nil, errors.Wrap(err, "generate client secret")
	}
	saltedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "encrypt client secret")
	}
	client.Secret = string(saltedSecret)

	err = u.clientsRepo.InsertClient(ctx, *client)
	if err != nil {
		return nil, errors.Wrap(err, "insert client")
	}

	client.Secret = secret
	return client, nil
}
