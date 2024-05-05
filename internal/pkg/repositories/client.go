package repositories

import (
	"context"

	"authorization-server/internal/pkg/domain"
)

type Storage interface {
	Insert(ctx context.Context, key string, value interface{}) error
}

type ClientsRepository struct {
	storage Storage
}

func NewClientsRepository(storage Storage) *ClientsRepository {
	return &ClientsRepository{
		storage: storage,
	}
}

func (r *ClientsRepository) InsertClient(ctx context.Context, client domain.Client) error {
	err := r.storage.Insert(ctx, client.ID, &client)
	if err != nil {
		return err
	}

	return nil
}
