package repositories

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"oauth2/internal/generated/authorization-server/public/model"
	"oauth2/internal/generated/authorization-server/public/table"

	"github.com/jmoiron/sqlx"
	"oauth2/internal/pkg/domain"
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

	_, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientsRepository) GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	stmt := table.Clients.
		SELECT(table.Clients.AllColumns).
		WHERE(table.Clients.ID.EQ(postgres.UUID(id)))

	var client model.Clients
	err := stmt.QueryContext(ctx, r.db, &client)
	if err != nil {
		return nil, err
	}

	return mapClientToDomain(client), nil
}

func mapClientToModel(client domain.Client) model.Clients {
	return model.Clients{
		ID:             client.ID,
		Name:           client.Name,
		URL:            client.URL,
		RedirectURI:    client.RedirectURI,
		Secret:         client.Secret,
		IsConfidential: client.IsConfidential,
	}
}

func mapClientToDomain(client model.Clients) *domain.Client {
	return &domain.Client{
		ID:             client.ID,
		Name:           client.Name,
		URL:            client.URL,
		RedirectURI:    client.RedirectURI,
		Secret:         client.Secret,
		IsConfidential: client.IsConfidential,
	}
}
