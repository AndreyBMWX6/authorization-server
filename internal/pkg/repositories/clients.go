package repositories

import (
	"context"

	"authorization-server/internal/generated/authorization-server/public/model"
	"authorization-server/internal/generated/authorization-server/public/table"

	"authorization-server/internal/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type ClientsRepository struct {
	db *sqlx.DB
}

func NewClientsRepository(db *sqlx.DB) *ClientsRepository {
	return &ClientsRepository{
		db: db,
	}
}

func (r *ClientsRepository) InsertClient(ctx context.Context, client domain.Client) error {
	m := mapClientToModel(client)

	stmt := table.Clients.
		INSERT(table.Clients.AllColumns).
		MODEL(m)

	_, err := stmt.Exec(r.db)
	if err != nil {
		return err
	}

	return nil
}

func mapClientToModel(client domain.Client) model.Clients {
	return model.Clients{
		ID:          client.ID,
		Name:        client.Name,
		URL:         client.URL,
		RedirectURI: client.RedirectURI,
		Secret:      client.Secret,
	}
}
